package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ab-testing-service/internal/models"
	"github.com/ab-testing-service/internal/proxy"
)

func (s *Storage) SaveProxyConfig(ctx context.Context, cfg proxy.Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return s.Redis.Set(ctx, "proxy:"+cfg.ID, data, 0).Err()
}

func (s *Storage) LoadProxyConfigs(ctx context.Context) ([]proxy.Config, error) {
	keys, err := s.Redis.Keys(ctx, "proxy:*").Result()
	if err != nil {
		return nil, err
	}

	var configs []proxy.Config
	for _, key := range keys {
		data, err := s.Redis.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}

		var cfg proxy.Config
		if err := json.Unmarshal(data, &cfg); err != nil {
			continue
		}
		configs = append(configs, cfg)
	}

	return configs, nil
}

func (s *Storage) InvalidateProxyCache(ctx context.Context, proxyID string) error {
	return s.Redis.Del(ctx, fmt.Sprintf("proxy:%s", proxyID)).Err()
}

func (s *Storage) GetProxyConfig(ctx context.Context, proxyID string) (proxy.Config, error) {
	data, err := s.Redis.Get(ctx, fmt.Sprintf("proxy:%s", proxyID)).Bytes()
	if err != nil {
		return proxy.Config{}, err
	}

	var cfg proxy.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return proxy.Config{}, err
	}
	return cfg, nil
}

func (s *Storage) GetProxy(ctx context.Context, id string) (*models.Proxy, error) {
	// Try Redis first
	key := fmt.Sprintf("proxy:%s", id)
	data, err := s.Redis.Get(ctx, key).Bytes()
	if err == nil {
		var proxy models.Proxy
		if err := json.Unmarshal(data, &proxy); err == nil {
			return &proxy, nil
		}
	}

	// Fallback to PostgreSQL
	var conditionJSON models.RouteCondition

	p, err := s.q.GetProxy(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(p.Condition) > 0 {
		if err := json.Unmarshal(p.Condition, &conditionJSON); err != nil {
			return nil, fmt.Errorf("failed to unmarshal condition: %w", err)
		}
	}

	proxyModel := models.Proxy{
		ID:        p.ID,
		ListenURL: p.ListenUrl,
		Mode:      models.ProxyMode(p.Mode),
		Condition: &conditionJSON,
		PathKey:   p.PathKey,
		CreatedAt: p.CreatedAt.Time,
		UpdatedAt: p.UpdatedAt.Time,
	}

	targets, err := s.q.GetTargetsByProxyID(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, target := range targets {
		target := models.Target{
			ProxyID:  p.ID,
			ID:       target.ID,
			URL:      target.Url,
			Weight:   target.Weight,
			IsActive: target.IsActive,
		}
		proxyModel.Targets = append(proxyModel.Targets, target)
	}

	// Cache in Redis
	if data, err := json.Marshal(proxyModel); err == nil {
		s.Redis.Set(ctx, key, data, proxyTTL)
	}

	return &proxyModel, nil
}

func (s *Storage) GetTargets(ctx context.Context, proxyID string) ([]*models.Target, error) {
	// Try Redis first
	key := fmt.Sprintf("targets:%s", proxyID)

	data, err := s.Redis.Get(ctx, key).Bytes()
	if err == nil {
		var targets []*models.Target
		if err := json.Unmarshal(data, &targets); err == nil {
			return targets, nil
		}
	}

	// Fallback to PostgreSQL
	rows, err := s.q.GetTargetsByProxyID(ctx, proxyID)
	if err != nil {
		return nil, err
	}

	targets := make([]*models.Target, 0, len(rows))

	for _, item := range rows {
		targets = append(targets, &models.Target{
			ID:       item.ID,
			ProxyID:  proxyID,
			URL:      item.Url,
			Weight:   item.Weight,
			IsActive: item.IsActive,
		})
	}

	// Cache in Redis
	if data, err := json.Marshal(targets); err == nil {
		s.Redis.Set(ctx, key, data, targetTTL)
	}

	return targets, nil
}
