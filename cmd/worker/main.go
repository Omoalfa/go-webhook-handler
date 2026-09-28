package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Omoalfa/event-platform/internal/config"
	"github.com/Omoalfa/event-platform/internal/database"
	"github.com/Omoalfa/event-platform/internal/events"
	"github.com/Omoalfa/event-platform/internal/worker"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
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
	defer redisClient.Close()

	eventRepo := events.NewPGRepository(dbPool)
	processorService := events.NewProcessor(eventRepo)

	eventWorker, err := worker.NewWorker(cfg.RedisURL)
	if err != nil {
		log.Fatal(err)
	}

	processor := worker.NewEventProcessor(processorService)

	log.Println("Starting event worker...")

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- eventWorker.Start(processor)
	}()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(
		shutdownSignal,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	select {
	case err := <-serverErrors:
		if err != nil {
			log.Fatalf("Worker failed: %v", err)
		}

	case sig := <-shutdownSignal:
		log.Printf("Received shutdown signal %s", sig)
		eventWorker.Shutdown()
	}

}
