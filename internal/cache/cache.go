package cache

import (
	"context"
	"crypto/internal/models"
	"crypto/internal/repository"
	"strconv"
	"sync"
	"time"
	"unsafe"

	"github.com/pkg/errors"
)

type WrapWallet struct {
	wallet *models.Address
	expiry time.Time
	size   int64
}

type CacheDecorator struct {
	WalletRepo repository.WalletProvider
	ttl        time.Duration

	mu      sync.RWMutex
	Wallets map[uint64]WrapWallet
}

func CacheNewDecorator(repo repository.WalletProvider, ttl time.Duration) *CacheDecorator {
	return &CacheDecorator{
		WalletRepo: repo,
		Wallets:    make(map[uint64]WrapWallet),
		ttl:        ttl,
	}
}

func (c *CacheDecorator) Get(id uint64) (*models.Address, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	wallet, ok := c.Wallets[id]

	return wallet.wallet, ok
}

func (c *CacheDecorator) Set(wallet *models.Address) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	walletSize := c.sizeCache()

	c.Wallets[wallet.ID] = WrapWallet{
		wallet: wallet,
		expiry: time.Now().Add(c.ttl),
		size:   walletSize,
	}
}

func (c *CacheDecorator) CreateAddress(ctx context.Context, req *models.AddressRequest) (*models.Address, error) {
	wallet, err := c.WalletRepo.CreateAddress(ctx, req)
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

	wallet, err := c.WalletRepo.GetID(ctx, id)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to get wallet by id")
	}

	c.Set(wallet)
	return wallet, nil
}

func (c *CacheDecorator) GetAllWallets(ctx context.Context) ([]models.Address, error) {
	return c.WalletRepo.GetAllWallets(ctx)
}

func (c *CacheDecorator) EditTag(ctx context.Context, req *models.TagUpdateRequest) error {
	if err := c.WalletRepo.EditTag(ctx, req); err != nil {
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
	for id, w := range c.Wallets {
		if now.After(w.expiry) {
			delete(c.Wallets, id)
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
	return len(c.Wallets)
}

func (c *CacheDecorator) MemoryUsage() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var totalsize int64
	for _, w := range c.Wallets {
		totalsize += w.size
	}
	return totalsize
}

func (c *CacheDecorator) copy(copyTtl time.Duration, copyWalletRepo repository.WalletProvider) *CacheDecorator {
	copyWallets := make(map[uint64]WrapWallet)
	for k, v := range c.Wallets {
		copyWallets[k] = WrapWallet{
			wallet: &models.Address{
				ID:            v.wallet.ID,
				WalletAddress: v.wallet.WalletAddress,
				ChainName:     v.wallet.ChainName,
				CryptoName:    v.wallet.CryptoName,
				Tag:           v.wallet.Tag,
				Balance:       v.wallet.Balance,
			},
			expiry: v.expiry,
			size:   v.size,
		}
	}

	return &CacheDecorator{
		WalletRepo: copyWalletRepo,
		ttl:        copyTtl,
		Wallets:    copyWallets,
	}
}

func (c *CacheDecorator) sizeCache() int64 {
	replica := c.copy(c.ttl, c.WalletRepo)

	var size int64

	size += int64(unsafe.Sizeof(*replica))

	size += int64(unsafe.Sizeof(replica.Wallets))
	for k, v := range replica.Wallets {
		size += int64(int64(unsafe.Sizeof(k)) + int64(unsafe.Sizeof(v)))
		if v.wallet != nil {
			size += int64(unsafe.Sizeof(v.wallet))
			size += int64(len(strconv.FormatUint(v.wallet.ID, 10)))
			size += int64(len(v.wallet.WalletAddress))
			size += int64(len(v.wallet.ChainName))
			size += int64(len(v.wallet.CryptoName))
			size += int64(len(v.wallet.Tag))
			size += int64(len(strconv.FormatInt(v.wallet.Balance, 10)))
		}
	}
	return size
}
