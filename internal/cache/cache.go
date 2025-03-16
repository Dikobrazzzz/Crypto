package cache

import (
	"context"
	"crypto/internal/models"
	"crypto/internal/repository"
	"sync"
	"time"
)

type wrapWallet struct {
	wallet *models.Address
	expiry time.Time
}

type CacheDecorator struct {
	walletRepo       repository.WalletProvider
	wallets          map[uint64]wrapWallet
	allWallets       []models.Address
	allWalletsCached bool
	mu               sync.Mutex
	ttl              time.Duration
}

func CacheNewDecorator(repo repository.WalletProvider) *CacheDecorator {
	c := &CacheDecorator{
		walletRepo:       repo,
		wallets:          make(map[uint64]wrapWallet),
		allWallets:       []models.Address{},
		allWalletsCached: false,
		ttl:              5 * time.Minute,
	}
	return c
}

func (c *CacheDecorator) Get(id uint64) (*models.Address, bool) {
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
		return nil, err
	}

	c.Set(wallet)

	return wallet, err
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
	err := c.walletRepo.EditTag(ctx, req)
	if err != nil {
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
