package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"wb-tech-L0/domain/model"

	"github.com/go-playground/validator/v10"
	"github.com/segmentio/kafka-go"
)

type OrderService interface {
	SaveOrder(ctx context.Context, order *model.Order) error
}

func HandleMessage(ctx context.Context, service OrderService) MessageHandler {
	return func(v *validator.Validate, msg kafka.Message) error {
		var order model.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			fmt.Println("Couldn't unmarshal message")
			return err
		}

		if err := v.Struct(order); err != nil {
			return err
		}

		return service.SaveOrder(ctx, &order)
	}
}
