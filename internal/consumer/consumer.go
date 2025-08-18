package consumer

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type MessageHandler func(kafka.Message) error

type Consumer struct {
	reader  *kafka.Reader
	handler MessageHandler
}

func New(cfg Config, handler MessageHandler) *Consumer {
	return &Consumer{reader: kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		Topic:   cfg.Topic,
		GroupID: cfg.GroupID,
	}),
		handler: handler}
}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Println("Consumer is stopped")
			return nil
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return nil
				}
				log.Println("Error while reading message")
				time.Sleep(2 * time.Second)
				continue
			}

			if err := c.handler(msg); err != nil {
				log.Println("Error while handling message.", err)
			} else {
				if err := c.reader.CommitMessages(ctx); err != nil {
					log.Println("Error while commiting message")
				}
			}
		}

	}
}
func (c *Consumer) Close() error {
	return c.reader.Close()
}
