package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/lib/pq"
	_ "github.com/lib/pq"

	"github.com/ab-testing-service/internal/models"
	"github.com/ab-testing-service/internal/proxy"
)

type Storage struct {
	db    *sql.DB
	Redis *redis.Client
}

const (
	proxyTTL  = 1 * time.Hour
	targetTTL = 1 * time.Hour
)

type Tx struct {
	tx *sql.Tx
}

func NewStorage(db *sql.DB, redis *redis.Client) *Storage {
	return &Storage{
		db:    db,
		Redis: redis,
	}
}

func (s *Storage) BeginTx(ctx context.Context) (*Tx, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Tx{tx: tx}, nil
}

func (tx *Tx) Commit() error {
	return tx.tx.Commit()
}

func (tx *Tx) Rollback() error {
	return tx.tx.Rollback()
}

func (tx *Tx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return tx.tx.ExecContext(ctx, query, args...)
}

func (s *Storage) UpdateProxyCondition(ctx context.Context, proxyID string, condition *models.RouteCondition) error {
	conditionJSON, err := json.Marshal(condition)
	if err != nil {
		return fmt.Errorf("failed to marshal condition: %w", err)
	}

	_, err = s.db.ExecContext(ctx,
		`UPDATE proxies SET condition = $1, updated_at = $2 WHERE id = $3`,
		conditionJSON, time.Now(), proxyID,
	)
	return err
}

func (s *Storage) UpdateProxyConditionWithTx(ctx context.Context, tx *Tx, proxyID string, condition *models.RouteCondition) error {
	conditionJSON, err := json.Marshal(condition)
	if err != nil {
		return fmt.Errorf("failed to marshal condition: %w", err)
	}

	_, err = tx.tx.ExecContext(ctx,
		`UPDATE proxies SET condition = $1, updated_at = $2 WHERE id = $3`,
		conditionJSON, time.Now(), proxyID,
	)
	return err
}

func (s *Storage) SaveVisit(ctx context.Context, visit *models.Visit) error {
	visit.ID = uuid.New().String()
	visit.CreatedAt = time.Now()

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO visits (id, proxy_id, target_id, user_id, rid, rrid, ruid, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		visit.ID, visit.ProxyID, visit.TargetID, visit.UserID,
		visit.RID, visit.RRID, visit.RUID, visit.CreatedAt,
	)
	return err
}

func (s *Storage) UpdateTargets(ctx context.Context, proxyID string, targets []models.Target) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete existing targets
	_, err = tx.ExecContext(ctx, `DELETE FROM targets WHERE proxy_id = $1`, proxyID)
	if err != nil {
		return err
	}

	// Insert new targets
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO targets (id, proxy_id, url, weight, is_active)
		VALUES ($1, $2, $3, $4, $5)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, target := range targets {
		_, err = stmt.ExecContext(ctx,
			target.ID,
			proxyID,
			target.URL,
			target.Weight,
			target.IsActive,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Storage) UpdateTargetsWithTx(ctx context.Context, tx *Tx, proxyID string, targets []models.Target) error {
	// Delete existing targets
	_, err := tx.tx.ExecContext(ctx,
		`DELETE FROM targets WHERE proxy_id = $1`,
		proxyID,
	)
	if err != nil {
		return err
	}

	// Insert new targets
	for _, target := range targets {
		_, err = tx.tx.ExecContext(ctx,
			`INSERT INTO targets (id, proxy_id, url, weight, is_active)
			VALUES ($1, $2, $3, $4, $5)`,
			target.ID, proxyID, target.URL, target.Weight, target.IsActive,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) GetProxies(ctx context.Context) ([]proxy.Config, error) {
	var proxies []proxy.Config
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, listen_url, mode, condition, tags
		FROM proxies ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query proxies: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Proxy
		var conditionJSON []byte
		if err := rows.Scan(&p.ID, &p.ListenURL, &p.Mode, &conditionJSON, pq.Array(&p.Tags)); err != nil {
			return nil, fmt.Errorf("failed to scan proxy: %w", err)
		}
		if len(conditionJSON) > 0 {
			p.Condition = &models.RouteCondition{}
			if err := json.Unmarshal(conditionJSON, p.Condition); err != nil {
				return nil, fmt.Errorf("failed to unmarshal condition: %w", err)
			}
		}
		proxies = append(proxies, proxy.Config{
			ID:        p.ID,
			ListenURL: p.ListenURL,
			Mode:      models.ProxyMode(p.Mode),
			Condition: (*proxy.Condition)(p.Condition),
			Tags:      p.Tags,
		})
	}
	return proxies, nil
}
