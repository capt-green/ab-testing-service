package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
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

	dbpool, err := pgxpool.New(ctx, fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	))
	defer dbpool.Close()

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Kafka writer
	kafkaURL := cfg.Kafka.KafkaURL
	topic := cfg.Kafka.Topic
	if err := createTopic(topic, kafkaURL); err != nil {
		log.Fatal("Failed to create topic:", err)
	}
	kw := getKafkaWriter(kafkaURL, topic)
	defer kw.Close()

	// Initialize storage
	store := storage.NewStorage(dbpool, rdb)

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

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func createTopic(topic, kafkaURL string) error {
	conn, err := kafka.Dial("tcp", kafkaURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		return err
	}

	return nil
}
