package testhelpers

import (
	"time"
	"wb-tech-L0/domain/model"

	"github.com/google/uuid"
)

func NewTestOrder() *model.Order {
	return &model.Order{
		OrderUID:          uuid.New(),
		TrackNumber:       "",
		Entry:             "",
		Delivery:          model.Delivery{},
		Payment:           model.Payment{},
		Items:             nil,
		Locale:            "",
		InternalSignature: "",
		CustomerID:        "",
		DeliveryService:   "",
		ShardKey:          "",
		SmID:              0,
		DateCreated:       time.Time{},
		OofShard:          "",
	}
}
