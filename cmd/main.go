package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wb-tech-L0/internal/cache"
	"wb-tech-L0/internal/consumer"
	"wb-tech-L0/internal/db"
	"wb-tech-L0/internal/httpapi"
	"wb-tech-L0/internal/repository"
	"wb-tech-L0/internal/service"

	"golang.org/x/sync/errgroup"
)

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	g, gctx := errgroup.WithContext(ctx)

	postgresDb, err := db.New(gctx, dbURL)
	if err != nil {
		log.Fatalf("Couldn't connect to database")
		return
	}
	defer postgresDb.Close()

	redisAddress := os.Getenv("REDIS_ADDRESS")
	redisConfig := cache.Config{Addr: redisAddress}
	redisCache := cache.NewRedisCache(redisConfig)
	defer func(redisCache *cache.RedisCache) {
		if err := redisCache.Close(); err != nil {
			log.Println("Error closing cache")
		}
	}(redisCache)

	repo := repository.NewCachedDB(postgresDb, redisCache)
	svc := service.New(gctx, repo)

	consumerConfig := consumer.Config{
		Brokers: []string{"kafka:9092"},
		Topic:   "topic",
		GroupID: "group",
	}
	cons := consumer.New(
		consumerConfig,
		consumer.HandleMessage(gctx, svc),
	)

	hndlr := httpapi.NewHandler(svc)
	router := httpapi.NewRouter(hndlr)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	g.Go(func() error {
		return server.ListenAndServe()
	})

	g.Go(func() error {
		return cons.Run(gctx)
	})

	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = server.Shutdown(shutdownCtx)
	if err != nil {
		log.Printf("http shutdown error: %v", err)
	}
	err = cons.Close()
	if err != nil {
		log.Printf("consumer shutdown error: %v", err)
	}

	if err = g.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("exit with error: %v", err)
	}
}
