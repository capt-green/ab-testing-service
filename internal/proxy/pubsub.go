package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

const (
	proxySettingsChannel = "proxy:settings:changes"
)

type SettingsChangeMessage struct {
	ProxyID   string `json:"proxy_id"`
	SenderID  string `json:"sender_id"` // Unique ID of the service instance that sent the message
	Operation string `json:"operation"` // e.g., "update", "delete"
}

type RedisPubSub struct {
	client     *redis.Client
	instanceID string
	onUpdate   func(ctx context.Context, proxyID string) error // Callback function for handling updates
}

func NewRedisPubSub(redisClient *redis.Client, updateCallback func(ctx context.Context, proxyID string) error) *RedisPubSub {
	return &RedisPubSub{
		client:     redisClient,
		instanceID: uuid.New().String(),
		onUpdate:   updateCallback,
	}
}

// StartSubscriber starts listening for proxy settings changes
func (ps *RedisPubSub) StartSubscriber(ctx context.Context) error {
	pubsub := ps.client.Subscribe(ctx, proxySettingsChannel)
	defer pubsub.Close()

	ch := pubsub.Channel()

	go func() {
		for msg := range ch {
			var changeMsg SettingsChangeMessage
			if err := json.Unmarshal([]byte(msg.Payload), &changeMsg); err != nil {
				log.Printf("Error unmarshaling settings change message: %v", err)
				continue
			}

			// Skip if we're the sender
			if changeMsg.SenderID == ps.instanceID {
				continue
			}

			// Call the update callback to handle the change
			if err := ps.onUpdate(ctx, changeMsg.ProxyID); err != nil {
				log.Printf("Error handling settings change for proxy %s: %v", changeMsg.ProxyID, err)
			}
		}
	}()

	return nil
}

// PublishSettingsChange notifies other instances about proxy settings changes
func (ps *RedisPubSub) PublishSettingsChange(ctx context.Context, proxyID string) error {
	msg := SettingsChangeMessage{
		ProxyID:   proxyID,
		SenderID:  ps.instanceID,
		Operation: "update",
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error marshaling settings change message: %w", err)
	}

	return ps.client.Publish(ctx, proxySettingsChannel, payload).Err()
}
