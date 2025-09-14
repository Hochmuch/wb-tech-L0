package model

import (
	"github.com/google/uuid"
)

type Payment struct {
	Transaction  uuid.UUID `json:"transaction" validate:"required,uuid4"`
	RequestID    string    `json:"request_id"`
	Currency     string    `json:"currency" validate:"required,iso4217"`
	Provider     string    `json:"provider" validate:"required"`
	Amount       int       `json:"amount" validate:"gte=1"`
	PaymentDT    int64     `json:"payment_dt" validate:"required"`
	Bank         string    `json:"bank" validate:"required"`
	DeliveryCost int       `json:"delivery_cost" validate:"gte=0"`
	GoodsTotal   int       `json:"goods_total" validate:"gte=1"`
	CustomFee    int       `json:"custom_fee" validate:"gte=0"`
}
