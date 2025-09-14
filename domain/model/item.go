package model

import "github.com/google/uuid"

type Item struct {
	ID          uuid.UUID `json:"-"`
	OrderUID    uuid.UUID `json:"-"`
	ChrtID      int       `json:"chrt_id" validate:"required"`
	TrackNumber string    `json:"track_number" validate:"required"`
	Price       int       `json:"price" validate:"gte=1"`
	Rid         string    `json:"rid" validate:"required,uuid4"`
	Name        string    `json:"name" validate:"required"`
	Sale        int       `json:"sale" validate:"gte=0"`
	Size        string    `json:"size" validate:"required,numeric"`
	TotalPrice  int       `json:"total_price" validate:"gte=1"`
	NmID        int       `json:"nm_id" validate:"required"`
	Brand       string    `json:"brand" validate:"required"`
	Status      int       `json:"status" validate:"required"`
}
