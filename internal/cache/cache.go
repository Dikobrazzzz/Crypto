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
	size   int64
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
		ttl:              ttl,
	}
}

func (c *CacheDecorator) Get(id uint64) (*models.Address, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	wallet, ok := c.wallets[id]

	return wallet.wallet, ok
}

func (c *CacheDecorator) Set(wallet *models.Address) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	walletSize := SizeCache(wallet)

	c.wallets[wallet.ID] = wrapWallet{
		wallet: wallet,
		expiry: time.Now().Add(c.ttl),
		size:   walletSize,
	}
}

func (c *CacheDecorator) CreateAddress(ctx context.Context, req *models.AddressRequest) (*models.Address, error) {
	wallet, err := c.walletRepo.CreateAddress(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "The operation cannot be performed")
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
		return nil, errors.Wrap(err, "Failed to get wallet by id")
	}

	c.Set(wallet)
	return wallet, nil
}

func (c *CacheDecorator) GetAllWallets(ctx context.Context) ([]models.Address, error) {
	return c.walletRepo.GetAllWallets(ctx)
}

func (c *CacheDecorator) EditTag(ctx context.Context, req *models.TagUpdateRequest) error {
	if err := c.walletRepo.EditTag(ctx, req); err != nil {
		return errors.Wrap(err, "Failed to edit tag")
	}

	if cachedAddr, ok := c.Get(req.ID); ok {
		updatedAddr := *cachedAddr
		updatedAddr.Tag = req.Tag
		c.Set(&updatedAddr)
	}
	return nil
}

func (c *CacheDecorator) removeExpired() {
	c.mu.RLock()
	defer c.mu.RUnlock()

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

func (c *CacheDecorator) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.wallets)
}

func (c *CacheDecorator) MemoryUsage() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var totalsize int64
	for _, w := range c.wallets {
		totalsize += w.size
	}
	return totalsize
}

func SizeCache(w *models.Address) int64 {
	var size int64

	size += 4
	size += int64(len(w.WalletAddress))
	size += int64(len(w.ChainName))
	size += int64(len(w.CryptoName))
	size += int64(len(w.Tag))
	size += 16
	return size
}
