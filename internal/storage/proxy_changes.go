package storage

import (
	"context"
	"fmt"

	"github.com/ab-testing-service/internal/models"
)

func (s *Storage) GetProxyChanges(ctx context.Context, proxyID string, limit, offset int) ([]models.ProxyChange, error) {
	rows, err := s.q.GetProxyChangesByProxyID(ctx, &GetProxyChangesByProxyIDParams{
		ProxyID: proxyID,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to query proxy changes: %w", err)
	}

	var changes []models.ProxyChange

	for _, change := range rows {
		changes = append(changes, models.ProxyChange{
			ID:            change.ID,
			ProxyID:       change.ProxyID,
			ChangeType:    models.ChangeType(change.ChangeType),
			PreviousState: change.PreviousState,
			NewState:      change.NewState,
			CreatedAt:     change.CreatedAt.Time,
			CreatedBy:     change.CreatedBy,
		})
	}

	return changes, nil
}
