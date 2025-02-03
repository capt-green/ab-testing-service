package proxy

import (
	"log"
	"sync"
	"time"
)

type TargetStats struct {
	RequestCount int64
	ErrorCount   int64
	LastUpdated  time.Time
	UniqueUsers  map[string]struct{} // Track unique users by ID/IP
}

type Stats struct {
	mu      sync.RWMutex
	ProxyID string
	Targets map[string]*TargetStats // key is target ID
}

func NewProxyStats(proxyID string) *Stats {
	return &Stats{
		ProxyID: proxyID,
		Targets: make(map[string]*TargetStats),
	}
}

func (s *Stats) IncrementRequestsWithUser(targetID string, userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Targets[targetID]; !exists {
		s.Targets[targetID] = &TargetStats{
			UniqueUsers: make(map[string]struct{}),
		}
	}
	s.Targets[targetID].RequestCount++
	s.Targets[targetID].UniqueUsers[userID] = struct{}{}
	s.Targets[targetID].LastUpdated = time.Now()
	log.Printf("Request count for target %s: %d", targetID, s.Targets[targetID].RequestCount)
}

func (s *Stats) IncrementErrors(targetID string, userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Targets[targetID]; !exists {
		s.Targets[targetID] = &TargetStats{
			UniqueUsers: make(map[string]struct{}),
		}
	}
	s.Targets[targetID].ErrorCount++
	s.Targets[targetID].UniqueUsers[userID] = struct{}{}
	s.Targets[targetID].LastUpdated = time.Now()
}

func (s *Stats) GetStats() map[string]*TargetStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := make(map[string]*TargetStats)
	for id, target := range s.Targets {
		stats[id] = &TargetStats{
			RequestCount: target.RequestCount,
			ErrorCount:   target.ErrorCount,
			LastUpdated:  target.LastUpdated,
			UniqueUsers:  make(map[string]struct{}),
		}
		for user := range target.UniqueUsers {
			stats[id].UniqueUsers[user] = struct{}{}
		}
	}
	//log.Printf("Stats for proxy %s: %v", s.ProxyID, stats)
	return stats
}

func (s *Stats) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Targets = make(map[string]*TargetStats)
}
