package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Omoalfa/event-platform/internal/config"
	"github.com/Omoalfa/event-platform/internal/database"
	"github.com/Omoalfa/event-platform/internal/events"
	"github.com/Omoalfa/event-platform/internal/queue"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Initialize the application
	ctx := context.Background()
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	dbPool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	redisClient, err := database.NewRedisClient(ctx, cfg.RedisURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("REDIS pinged!")
	defer redisClient.Close()

	asynqQueue, err := queue.NewAsynqQueue(cfg.RedisURL)
	if err != nil {
		log.Fatal(err)
	}
	defer asynqQueue.Close()

	r := chi.NewRouter()

	eventRepo := events.NewPGRepository(dbPool)
	eventService := events.NewService(eventRepo, asynqQueue)
	eventHandler := events.NewEventsHandler(eventService, cfg.WebhookSecret)

	events.HandlerFromMux(eventHandler, r)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf(
			"Starting API Server on port %s in %s environment",
			cfg.Port,
			cfg.Environment,
		)

		serverErrors <- server.ListenAndServe()
	}()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(
		shutdownSignal,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}

	case sig := <-shutdownSignal:
		log.Printf("Received shutdown signal: %s", sig)

		ctx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("server shutdown failed: %v", err)
		}
	}
}
