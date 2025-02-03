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
			log.Printf("Error getting stats for proxy %s", instance.Proxy.Config.ID)
			continue
		}

		// Get current stats
		currentStats := stats.GetStats()

		for targetID, targetStats := range currentStats {
			// Преобразуем map[string]struct{} в []string для JSON
			uniqueUsers := make([]string, 0, len(targetStats.UniqueUsers))
			for userID := range targetStats.UniqueUsers {
				uniqueUsers = append(uniqueUsers, userID)
			}

			statsMsg := map[string]interface{}{
				"proxy_id":      instance.Proxy.Config.ID,
				"target_id":     targetID,
				"timestamp":     time.Now().Unix(),
				"request_count": targetStats.RequestCount,
				"error_count":   targetStats.ErrorCount,
				"unique_users":  uniqueUsers,
			}

			msgBytes, err := json.Marshal(statsMsg)
			log.Printf("Statistics message: %v", statsMsg)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				continue
			}

			err = s.kafkaWriter.WriteMessages(ctx,
				kafka.Message{
					Value: msgBytes,
				},
			)
			if err != nil {
				log.Printf("Error writing message: %v", err)
			}
		}

		// Reset stats after successful sending
		stats.Reset()
	}

	//log.Printf("Statistics collected")
}
