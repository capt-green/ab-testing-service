package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"

	"github.com/ab-testing-service/internal/config"
	"github.com/ab-testing-service/internal/server"
	"github.com/ab-testing-service/internal/storage"
	"github.com/ab-testing-service/internal/supervisor"
)

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		DB:   cfg.Redis.DB,
	})
	defer rdb.Close()

	// Initialize PostgreSQL
	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Kafka writer
	kw := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topic,
	})
	defer kw.Close()

	// Initialize storage
	store := storage.NewStorage(db, rdb)

	// Create supervisor
	sup := supervisor.NewSupervisor(supervisor.Config{
		Config:      cfg,
		Storage:     store,
		KafkaWriter: kw,
	})

	// Create and start HTTP server
	srv := server.NewServer(cfg, sup, store)
	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Start supervisor
	go sup.Start(ctx)

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
		os.Exit(1)
	}

	if err := sup.Shutdown(ctx); err != nil {
		log.Printf("Supervisor shutdown error: %v", err)
		os.Exit(1)
	}
}
