package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"

	"github.com/ab-testing-service/internal/models"
)

func (s *Storage) UpdateProxyTags(ctx context.Context, proxyID string, tags []string) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE proxies SET tags = $1, updated_at = NOW() WHERE id = $2`,
		pq.Array(tags), proxyID,
	)
	if err != nil {
		return fmt.Errorf("failed to update proxy tags: %w", err)
	}
	return nil
}

func (s *Storage) GetAllTags(ctx context.Context) ([]string, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT DISTINCT UNNEST(tags) FROM proxies WHERE tags IS NOT NULL ORDER BY 1`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query tags: %w", err)
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (s *Storage) GetTags(proxyID string) []string {
	var tags []string

	if err := s.db.QueryRow(
		`SELECT tags FROM proxies WHERE id = $1`,
		proxyID,
	).Scan(pq.Array(&tags)); err != nil {
		return nil
	}

	return tags
}

func (s *Storage) GetProxiesByTags(ctx context.Context, tags []string) ([]*models.Proxy, error) {
	query := `
		SELECT DISTINCT p.id, p.listen_url, p.mode, p.condition, p.tags, p.created_at, p.updated_at
		FROM proxies p
		WHERE tags @> $1
		ORDER BY p.created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, pq.Array(tags))
	if err != nil {
		return nil, fmt.Errorf("failed to query proxies by tags: %w", err)
	}
	defer rows.Close()

	var proxies []*models.Proxy
	for rows.Next() {
		var proxy models.Proxy
		var conditionJSON []byte
		if err := rows.Scan(
			&proxy.ID,
			&proxy.ListenURL,
			&proxy.Mode,
			&conditionJSON,
			pq.Array(&proxy.Tags),
			&proxy.CreatedAt,
			&proxy.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan proxy: %w", err)
		}

		if len(conditionJSON) > 0 {
			if err := json.Unmarshal(conditionJSON, &proxy.Condition); err != nil {
				return nil, fmt.Errorf("failed to unmarshal condition: %w", err)
			}
		}

		proxies = append(proxies, &proxy)
	}

	return proxies, nil
}
