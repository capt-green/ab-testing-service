package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

type ProxyStats struct {
	ProxyID      string   `json:"proxy_id"`
	TargetID     string   `json:"target_id"`
	Timestamp    int64    `json:"timestamp"`
	RequestCount int      `json:"request_count"`
	ErrorCount   int      `json:"error_count"`
	UniqueUsers  []string `json:"unique_users"`
}

func main() {
	// Context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a WaitGroup to wait for the consumer to finish
	var wg sync.WaitGroup

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	kafkaURL := os.Getenv("kafkaURL")
	topic := os.Getenv("topic")
	groupID := os.Getenv("groupID")

	// Kafka configuration
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaURL},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	// PostgreSQL connection
	connStr := "postgresql://abtest:abtest@postgres:5432/abtest?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	// Prepare insert statement
	stmt, err := db.Prepare(`
		INSERT INTO proxy_stats (
			proxy_id, 
			target_id, 
			timestamp, 
			request_count, 
			error_count, 
			unique_users
		) VALUES ($1, $2, $3, $4, $5, $6)
	`)
	if err != nil {
		log.Fatal("Error preparing statement:", err)
	}

	// Start the consumer in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer r.Close()
		defer db.Close()
		defer stmt.Close()

		for {
			select {
			case <-ctx.Done():
				log.Println("Shutting down consumer...")
				return
			default:
				// Create a context with timeout for each message
				msgCtx, msgCancel := context.WithTimeout(ctx, 5*time.Second)

				// Read message
				m, err := r.ReadMessage(msgCtx)
				msgCancel() // Cancel the message context immediately after read attempt

				if err != nil {
					if ctx.Err() != nil {
						// Context was cancelled, break the loop
						return
					}
					log.Println("Error reading message:", err)
					continue
				}

				// Parse message
				var stats ProxyStats
				if err := json.Unmarshal(m.Value, &stats); err != nil {
					log.Println("Error unmarshaling message:", err)
					continue
				}

				// Convert unique users to JSONB
				uniqueUsersJSON, err := json.Marshal(stats.UniqueUsers)
				if err != nil {
					log.Println("Error marshaling unique users:", err)
					continue
				}

				// Convert Unix timestamp to time.Time
				timestamp := time.Unix(stats.Timestamp, 0)

				// Insert into database with context
				_, err = stmt.ExecContext(ctx,
					stats.ProxyID,
					stats.TargetID,
					timestamp,
					stats.RequestCount,
					stats.ErrorCount,
					uniqueUsersJSON,
				)
				if err != nil {
					log.Println("Error inserting into database:", err)
					continue
				}

				log.Printf("Successfully processed message for proxy %s\n", stats.ProxyID)
			}
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Received shutdown signal. Starting graceful shutdown...")

	// Trigger shutdown
	cancel()

	// Wait for consumer to finish with timeout
	shutdownChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(shutdownChan)
	}()

	// Wait for graceful shutdown with timeout
	select {
	case <-shutdownChan:
		log.Println("Graceful shutdown completed")
	case <-time.After(30 * time.Second):
		log.Println("Shutdown timed out after 30 seconds")
	}
}
