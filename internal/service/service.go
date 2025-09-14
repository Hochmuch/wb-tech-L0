package service

import (
	"context"
	"wb-tech-L0/domain/model"

	"github.com/google/uuid"
)

type Repository interface {
	SaveOrder(context.Context, *model.Order) error
	GetOrder(context.Context, uuid.UUID) (*model.Order, error)
	RestoreCache(ctx context.Context) error
}

type Service struct {
	repository Repository
}

func New(ctx context.Context, repository Repository) *Service {
	service := &Service{repository: repository}
	err := service.RestoreCache(ctx)
	if err != nil {
		return nil
	}
	return service
}

func (s *Service) SaveOrder(ctx context.Context, order *model.Order) error {
	return s.repository.SaveOrder(ctx, order)
}

func (s *Service) GetOrder(ctx context.Context, orderUID uuid.UUID) (*model.Order, error) {
	return s.repository.GetOrder(ctx, orderUID)
}

func (s *Service) RestoreCache(ctx context.Context) error {
	return s.repository.RestoreCache(ctx)
}
