package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/ab-testing-service/internal/models"
)

func (s *Storage) RecordProxyChange(ctx context.Context, tx *Tx, proxyID string, changeType models.ChangeType, previousState, newState interface{}, userID *string) error {
	previousJSON, err := json.Marshal(previousState)
	if err != nil {
		return fmt.Errorf("failed to marshal previous state: %w", err)
	}

	newJSON, err := json.Marshal(newState)
	if err != nil {
		return fmt.Errorf("failed to marshal new state: %w", err)
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO proxy_changes (id, proxy_id, change_type, previous_state, new_state, created_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		uuid.New().String(), proxyID, changeType, previousJSON, newJSON, time.Now(), userID,
	)
	if err != nil {
		return fmt.Errorf("failed to insert proxy change: %w", err)
	}

	return nil
}

func (s *Storage) GetProxyChanges(ctx context.Context, proxyID string, limit, offset int) ([]models.ProxyChange, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, proxy_id, change_type, previous_state, new_state, created_at, created_by
		FROM proxy_changes
		WHERE proxy_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		proxyID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query proxy changes: %w", err)
	}
	defer rows.Close()

	var changes []models.ProxyChange
	for rows.Next() {
		var change models.ProxyChange
		err := rows.Scan(
			&change.ID,
			&change.ProxyID,
			&change.ChangeType,
			&change.PreviousState,
			&change.NewState,
			&change.CreatedAt,
			&change.CreatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan proxy change: %w", err)
		}
		changes = append(changes, change)
	}

	return changes, nil
}
