package supervisor

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func (s *Supervisor) collectStatistics(ctx context.Context) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, instance := range s.proxies {
		if !instance.Started {
			continue
		}

		stats := instance.Proxy.GetStats()
		if stats == nil {
			continue
		}

		// Get current stats
		currentStats := stats.GetStats()

		// Create stats message
		statsMsg := map[string]interface{}{
			"proxy_id":     instance.Proxy.Config.ID,
			"timestamp":    time.Now().Unix(),
			"target_stats": currentStats,
		}

		// Send stats to Kafka if configured
		if s.kafkaWriter != nil {
			statsJSON, err := json.Marshal(statsMsg)
			if err != nil {
				log.Printf("Error marshaling stats for proxy %s: %v", instance.Proxy.Config.ID, err)
				continue
			}

			msg := kafka.Message{
				Key:   []byte(instance.Proxy.Config.ID),
				Value: statsJSON,
			}
			if err := s.kafkaWriter.WriteMessages(ctx, msg); err != nil {
				log.Printf("Error sending stats for proxy %s: %v", instance.Proxy.Config.ID, err)
			}
		}

		// Reset stats after successful sending
		stats.Reset()
	}
}
