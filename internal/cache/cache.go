package cache

import (
	"context"
	"crypto/internal/models"
	"crypto/internal/repository"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type wrapWallet struct {
	wallet *models.Address
	expiry time.Time
}

type CacheDecorator struct {
	walletRepo       repository.WalletProvider
	allWallets       []models.Address
	allWalletsCached bool

	mu      sync.RWMutex
	wallets map[uint64]wrapWallet
	ttl     time.Duration
}

func CacheNewDecorator(repo repository.WalletProvider, ttl time.Duration) *CacheDecorator {
	return &CacheDecorator{
		walletRepo:       repo,
		wallets:          make(map[uint64]wrapWallet),
		allWallets:       []models.Address{},
		allWalletsCached: false,
		ttl:              5 * time.Minute,
	}
}

func (c *CacheDecorator) Get(id uint64) (*models.Address, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	wallet, ok := c.wallets[id]

	return wallet.wallet, ok
}

func (c *CacheDecorator) Set(wallet *models.Address) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.wallets[wallet.ID] = wrapWallet{
		wallet: wallet,
		expiry: time.Now().Add(c.ttl),
	}
}

func (c *CacheDecorator) CreateAddress(ctx context.Context, req *models.AddressRequest) (*models.Address, error) {
	wallet, err := c.walletRepo.CreateAddress(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "The operation cannot be preformed")
	}

	c.Set(wallet)

	return wallet, nil
}

func (c *CacheDecorator) GetID(ctx context.Context, id uint64) (*models.Address, error) {
	if wallet, ok := c.Get(id); ok {
		return wallet, nil
	}

	wallet, err := c.walletRepo.GetID(ctx, id)

	if err != nil {
		return nil, err
	}

	c.Set(wallet)
	return wallet, nil
}

func (c *CacheDecorator) GetAllWallets(ctx context.Context) ([]models.Address, error) {
	if c.allWalletsCached {
		return c.allWallets, nil
	}

	wallets, err := c.walletRepo.GetAllWallets(ctx)
	if err != nil {
		return nil, err
	}

	c.allWallets = wallets
	c.allWalletsCached = true

	return wallets, nil
}

func (c *CacheDecorator) EditTag(ctx context.Context, req *models.TagUpdateRequest) error {
	if err := c.walletRepo.EditTag(ctx, req); err != nil {
		return err
	}

	if cachedAddr, ok := c.Get(req.ID); ok {
		updatedAddr := *cachedAddr
		updatedAddr.Tag = req.Tag
		c.Set(&updatedAddr)
	}
	return nil
}

func (c *CacheDecorator) removeExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for id, w := range c.wallets {
		if now.After(w.expiry) {
			delete(c.wallets, id)
		}
	}
}

func (c *CacheDecorator) StartCleaner(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Minute)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.removeExpired()
			case <-ctx.Done():
				return
			}
		}
	}()
}

// текущее количество элементов в кеше.
func (c *CacheDecorator) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.wallets)
}

// «вес» кеша.
// В данном упрощённом примере — просто Size() * 256 байт
func (c *CacheDecorator) MemoryUsage() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	const approximateBytesPerWallet = 256
	return int64(len(c.wallets)) * approximateBytesPerWallet
}
