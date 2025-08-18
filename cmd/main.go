package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"time"
	"wb-tech-L0/internal/cache"
	"wb-tech-L0/internal/consumer"
	"wb-tech-L0/internal/httpapi"
	"wb-tech-L0/internal/repository"
	"wb-tech-L0/internal/service"
)

func main() {
	databaseURL := "postgres://user:password@postgres:5432/db"
	time.Sleep(5 * time.Second)

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		fmt.Printf("Unable to parse config: %v\n", err)
		return
	}
	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		fmt.Printf("Unable to create connection pool: %v\n", err)
		return
	}
	defer pool.Close()

	repo := repository.New(ctx, pool)
	c := cache.New()
	serv := service.New(ctx, repo, c)
	cons := consumer.New(consumer.Config{Brokers: []string{"kafka:9092"},
		Topic:   "topic",
		GroupID: "group"},
		consumer.HandleMessage(ctx, serv))

	go func() {
		if err := cons.Run(ctx); err != nil {
			log.Fatalf("Kafka consumer error: %v", err)
		}
	}()

	hndlr := httpapi.NewHandler(serv)
	router := httpapi.NewRouter(hndlr)

	http.ListenAndServe(":8080", router)
}
