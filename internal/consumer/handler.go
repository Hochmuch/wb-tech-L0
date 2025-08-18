package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"wb-tech-L0/internal/model"
	"wb-tech-L0/internal/service"
)

func HandleMessage(ctx context.Context, service *service.Service) MessageHandler {
	return func(msg kafka.Message) error {
		var order model.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			fmt.Println("Couldn't unmarshal message")
			return err
		}
		return service.SaveOrder(ctx, &order)
	}
}
