package repository

import (
	"context"
	"errors"
	"wb-tech-L0/domain/model"

	"github.com/google/uuid"
)

var ErrOrderNotFound = errors.New("order not found")

type Repository interface {
	SaveOrder(context.Context, *model.Order) error
	GetOrder(context.Context, uuid.UUID) (*model.Order, error)
	GetDataForCache(ctx context.Context) ([]*model.Order, error)
	Close()
}

type Cache interface {
	AddOrder(context.Context, *model.Order) error
	GetOrder(context.Context, uuid.UUID) (*model.Order, error)
	Close() error
}

type CachedDB struct {
	db    Repository
	cache Cache
}

func NewCachedDB(repository Repository, cache Cache) *CachedDB {
	return &CachedDB{db: repository, cache: cache}
}

func (c *CachedDB) SaveOrder(ctx context.Context, order *model.Order) error {
	return c.db.SaveOrder(ctx, order)
}

func (c *CachedDB) GetOrder(ctx context.Context, orderUID uuid.UUID) (*model.Order, error) {
	order, err := c.cache.GetOrder(ctx, orderUID)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			order, err = c.db.GetOrder(ctx, orderUID)
			if err != nil {
				return nil, err
			}
			return order, nil
		}
		return nil, err
	}
	return order, nil
}

func (c *CachedDB) RestoreCache(ctx context.Context) error {
	dataToAddToCache, err := c.db.GetDataForCache(ctx)
	if err != nil {
		return err
	}
	for _, order := range dataToAddToCache {
		err = c.cache.AddOrder(ctx, order)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CachedDB) Close() error {
	c.db.Close()
	return c.cache.Close()
}
