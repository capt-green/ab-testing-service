package proxy

import (
	"sync"
	"time"
)

type TargetStats struct {
	RequestCount int64
	ErrorCount   int64
	LastUpdated  time.Time
}

type Stats struct {
	mu      sync.RWMutex
	Targets map[string]*TargetStats // key is target ID
}

func NewProxyStats() *Stats {
	return &Stats{
		Targets: make(map[string]*TargetStats),
	}
}

func (s *Stats) IncrementRequests(targetID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Targets[targetID]; !exists {
		s.Targets[targetID] = &TargetStats{}
	}
	s.Targets[targetID].RequestCount++
	s.Targets[targetID].LastUpdated = time.Now()
}

func (s *Stats) IncrementErrors(targetID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Targets[targetID]; !exists {
		s.Targets[targetID] = &TargetStats{}
	}
	s.Targets[targetID].ErrorCount++
	s.Targets[targetID].LastUpdated = time.Now()
}

func (s *Stats) GetStats() map[string]*TargetStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create a deep copy to avoid concurrent access issues
	stats := make(map[string]*TargetStats)
	for id, target := range s.Targets {
		stats[id] = &TargetStats{
			RequestCount: target.RequestCount,
			ErrorCount:   target.ErrorCount,
			LastUpdated:  target.LastUpdated,
		}
	}
	return stats
}

func (s *Stats) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Targets = make(map[string]*TargetStats)
}
