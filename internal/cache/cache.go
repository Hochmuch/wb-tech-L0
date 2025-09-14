package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
	"wb-tech-L0/domain/model"
	"wb-tech-L0/internal/repository"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	TTL    time.Duration
}

func NewRedisCache(cfg Config) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
	})

	return &RedisCache{
		client: client,
		TTL:    cfg.TTL,
	}
}

func (r *RedisCache) AddOrder(ctx context.Context, order *model.Order) error {
	orderJson, err := json.Marshal(order)
	if err != nil {
		log.Println(err)
		return err
	}
	err = r.client.Set(ctx, order.OrderUID.String(), string(orderJson), r.TTL).Err()
	return err
}

func (r *RedisCache) GetOrder(ctx context.Context, uid uuid.UUID) (*model.Order, error) {
	orderJson, err := r.client.Get(ctx, uid.String()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, repository.ErrOrderNotFound
		}
		return nil, fmt.Errorf("cache get order: %w", err)
	}

	var order model.Order
	err = json.Unmarshal([]byte(orderJson), &order)
	if err != nil {
		return nil, fmt.Errorf("cache unmarshal order: %w", err)
	}

	return &order, nil
}

func (r *RedisCache) Close() error {
	if r == nil || r.client == nil {
		return nil
	}
	return r.client.Close()
}
