package service_test

import (
	"context"
	"testing"
	"wb-tech-L0/internal/mocks"
	"wb-tech-L0/internal/service"
	"wb-tech-L0/internal/testhelpers"

	"github.com/golang/mock/gomock"
)

func TestService_GetOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		RestoreCache(ctx).
		Return(nil)

	svc := service.New(ctx, mockRepo)

	order := testhelpers.NewTestOrder()
	mockRepo.EXPECT().
		GetOrder(ctx, order.OrderUID).
		Return(order, nil)

	result, err := svc.GetOrder(ctx, order.OrderUID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != order {
		t.Fatalf("expected %v, got %v", order, result)
	}
}

func TestService_SaveOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		RestoreCache(ctx).
		Return(nil)

	svc := service.New(ctx, mockRepo)

	order := testhelpers.NewTestOrder()

	mockRepo.EXPECT().
		SaveOrder(ctx, order).
		Return(nil)

	err := svc.SaveOrder(ctx, order)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
