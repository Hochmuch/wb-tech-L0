package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"wb-tech-L0/internal/cache"
	"wb-tech-L0/internal/model"
	"wb-tech-L0/internal/repository"
)

type Service struct {
	repository *repository.Repository
	cache      *cache.Cache
}

func New(ctx context.Context, repository *repository.Repository, cache *cache.Cache) *Service {
	service := &Service{repository: repository, cache: cache}
	service.RestoreCache(ctx)
	return service
}

func (s *Service) SaveOrder(ctx context.Context, order *model.Order) error {
	if s.cache.HasSpace() {
		s.cache.AddOrder(order)
	}
	return s.repository.Save(ctx, order)
}

func (s *Service) GetOrder(ctx context.Context, orderUID uuid.UUID) (*model.Order, error) {
	if order, ok := s.cache.GetOrder(orderUID); ok {
		fmt.Println("RECORD FROM CACHE!!!")
		return order, nil
	}
	return s.repository.GetOrder(ctx, orderUID)
}

func (s *Service) RestoreCache(ctx context.Context) {
	orders, err := s.repository.GetDataForCache(ctx)
	if err != nil {
		return
	}
	s.cache.GetCache(orders)
}
