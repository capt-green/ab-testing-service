package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/ab-testing-service/internal/models"
)

func (s *Storage) CreateProxy(ctx context.Context, proxy *models.Proxy) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert proxy
	proxy.ID = uuid.New().String()
	proxy.CreatedAt = time.Now()
	proxy.UpdatedAt = proxy.CreatedAt

	var conditionJSON any = nil
	if proxy.Condition != nil {
		bytes, err := json.Marshal(proxy.Condition)
		if err != nil {
			return fmt.Errorf("failed to marshal condition: %w", err)
		}
		conditionJSON = bytes // Если есть данные, присваиваем []byte
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO proxies (id, listen_url, mode, condition, tags, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		proxy.ID, proxy.ListenURL, proxy.Mode, conditionJSON, pq.Array(proxy.Tags), proxy.CreatedAt, proxy.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert proxy: %w, condition: %s, condition type: %T, condition value: %+v", err, conditionJSON, proxy.Condition, proxy.Condition)
	}

	// Insert targets
	for i := range proxy.Targets {
		target := &proxy.Targets[i]
		target.ProxyID = proxy.ID

		_, err = tx.ExecContext(ctx,
			`INSERT INTO targets (id, proxy_id, url, weight, is_active)
			VALUES ($1, $2, $3, $4, $5)`,
			target.ID, target.ProxyID, target.URL, target.Weight, target.IsActive,
		)
		if err != nil {
			return fmt.Errorf("failed to insert target: %w", err)
		}
	}

	return tx.Commit()
}
