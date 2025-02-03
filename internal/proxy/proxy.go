package proxy

import (
	"fmt"
	"sync"

	"github.com/ab-testing-service/internal/models"
)

type Target struct {
	ID       string  `json:"id"`
	URL      string  `json:"url"`
	Weight   float64 `json:"weight"`
	IsActive bool    `json:"is_active"`
}

type Config struct {
	ID        string           `json:"id"`
	ListenURL string           `json:"listen_url"`
	Mode      models.ProxyMode `json:"mode"`
	PathKey   string           `json:"path_key,omitempty"`
	Targets   []Target         `json:"targets"`
	Condition *Condition       `json:"condition"`
	Tags      []string         `json:"tags"`
}

type Condition struct {
	Type      models.ConditionType `json:"type"`
	ParamName string               `json:"param_name"`
	Values    map[string]string    `json:"values"`
	Default   string               `json:"default"`
}

type Proxy struct {
	ID         string
	ListenURL  string
	Mode       models.ProxyMode
	PathKey    string
	Targets    []Target
	Config     Config
	mutex      sync.RWMutex
	metrics    *Metrics
	cookieName string
	stats      *Stats
}

func NewProxy(cfg Config) (*Proxy, error) {
	totalWeight, err := validate(cfg)
	if err != nil {
		return nil, err
	}

	// Normalize weights if total is not 1
	if totalWeight != 1 && totalWeight != 0 {
		for i := range cfg.Targets {
			cfg.Targets[i].Weight = cfg.Targets[i].Weight / totalWeight
		}
	}

	proxy := &Proxy{
		ID:         cfg.ID,
		ListenURL:  cfg.ListenURL,
		Mode:       cfg.Mode,
		Targets:    cfg.Targets,
		Config:     cfg,
		metrics:    newProxyMetrics(cfg.ID),
		cookieName: fmt.Sprintf("proxy_%s", cfg.ID),
		stats:      NewProxyStats(cfg.ID),
	}

	return proxy, nil
}

func validate(cfg Config) (float64, error) {
	if cfg.ID == "" {
		return 0, fmt.Errorf("proxy ID is required")
	}
	if cfg.ListenURL == "" {
		return 0, fmt.Errorf("listen URL is required")
	}
	if len(cfg.Targets) == 0 {
		return 0, fmt.Errorf("at least one target is required")
	}

	// Validate and normalize target weights
	var totalWeight float64
	for _, t := range cfg.Targets {
		if t.Weight < 0 {
			return 0, fmt.Errorf("target weight must be non-negative")
		}
		totalWeight += t.Weight
	}
	return totalWeight, nil
}

func (p *Proxy) UpdateTargets(targets []Target) {
	p.mutex.Lock()
	p.Targets = targets
	p.mutex.Unlock()
}

func (p *Proxy) GetStats() *Stats {
	return p.stats
}
