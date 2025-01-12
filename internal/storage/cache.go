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
	var proxy models.Proxy
	var conditionJSON []byte

	err = s.db.QueryRowContext(ctx,
		`SELECT id, listen_url, mode, condition, created_at, updated_at
		FROM proxies WHERE id = $1`,
		id,
	).Scan(&proxy.ID, &proxy.ListenURL, &proxy.Mode, &conditionJSON, &proxy.CreatedAt, &proxy.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if len(conditionJSON) > 0 {
		proxy.Condition = &models.RouteCondition{}
		if err := json.Unmarshal(conditionJSON, proxy.Condition); err != nil {
			return nil, fmt.Errorf("failed to unmarshal condition: %w", err)
		}
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT id, url, weight, is_active FROM targets WHERE proxy_id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var target models.Target
		target.ProxyID = proxy.ID
		if err := rows.Scan(&target.ID, &target.URL, &target.Weight, &target.IsActive); err != nil {
			return nil, err
		}
		proxy.Targets = append(proxy.Targets, target)
	}

	// Cache in Redis
	if data, err := json.Marshal(proxy); err == nil {
		s.Redis.Set(ctx, key, data, proxyTTL)
	}

	return &proxy, nil
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
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, proxy_id, url, weight, is_active FROM targets WHERE proxy_id = $1`,
		proxyID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []*models.Target
	for rows.Next() {
		var target models.Target
		if err := rows.Scan(&target.ID, &target.ProxyID, &target.URL, &target.Weight, &target.IsActive); err != nil {
			return nil, err
		}
		targets = append(targets, &target)
	}

	// Cache in Redis
	if data, err := json.Marshal(targets); err == nil {
		s.Redis.Set(ctx, key, data, targetTTL)
	}

	return targets, nil
}
