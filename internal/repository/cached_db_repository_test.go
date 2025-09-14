package repository_test

import (
	"context"
	"testing"
	"wb-tech-L0/domain/model"
	"wb-tech-L0/internal/mocks"
	"wb-tech-L0/internal/repository"
	"wb-tech-L0/internal/testhelpers"

	"github.com/golang/mock/gomock"
)

func TestCachedDB_GetOrder_CacheHasData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockDB := mocks.NewMockDB(ctrl)
	mockCache := mocks.NewMockCache(ctrl)

	order := testhelpers.NewTestOrder()

	mockCache.EXPECT().
		GetOrder(ctx, order.OrderUID).
		Return(order, nil)

	repo := repository.NewCachedDB(mockDB, mockCache)
	result, err := repo.GetOrder(ctx, order.OrderUID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != order {
		t.Fatalf("expected %v, got %v", order, result)
	}
}

func TestCachedDB_GetOrder_CacheHasNoData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockDB := mocks.NewMockDB(ctrl)
	mockCache := mocks.NewMockCache(ctrl)

	order := testhelpers.NewTestOrder()
	mockCache.EXPECT().
		GetOrder(ctx, order.OrderUID).
		Return(nil, repository.ErrOrderNotFound)
	mockDB.EXPECT().
		GetOrder(ctx, order.OrderUID).
		Return(order, nil)

	repo := repository.NewCachedDB(mockDB, mockCache)
	result, err := repo.GetOrder(ctx, order.OrderUID)

	if err != nil {
		t.Fatalf("unexpected error. %v", err)
	}

	if result != order {
		t.Fatalf("expected %v, got %v", order, result)
	}
}

func TestCachedDB_SaveOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	order := testhelpers.NewTestOrder()

	mockDB := mocks.NewMockDB(ctrl)
	mockCache := mocks.NewMockCache(ctrl)

	mockDB.EXPECT().
		SaveOrder(ctx, order).
		Return(nil)

	repo := repository.NewCachedDB(mockDB, mockCache)
	err := repo.SaveOrder(ctx, order)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCachedDB_RestoreCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	data := make([]*model.Order, 0, 10)
	for i := 0; i < 10; i++ {
		data = append(data, testhelpers.NewTestOrder())
	}

	mockDB := mocks.NewMockDB(ctrl)
	mockCache := mocks.NewMockCache(ctrl)

	mockDB.EXPECT().
		GetDataForCache(ctx).
		Return(data, nil)

	for _, order := range data {
		mockCache.EXPECT().
			AddOrder(ctx, order).
			Return(nil)
	}

	repo := repository.NewCachedDB(mockDB, mockCache)

	err := repo.RestoreCache(ctx)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Как я понимаю, такие тесты писать не нужно
func TestCachedDB_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockCache := mocks.NewMockCache(ctrl)

	mockDB.EXPECT().
		Close()

	mockCache.EXPECT().
		Close().
		Return(nil)

	repo := repository.NewCachedDB(mockDB, mockCache)

	err := repo.Close()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
