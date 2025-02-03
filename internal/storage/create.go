package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/ab-testing-service/internal/models"
)

func (s *Storage) CreateProxy(ctx context.Context, proxy *models.Proxy) error {
	err := pgx.BeginFunc(ctx, s.db, func(tx pgx.Tx) (err error) {
		repo := New(tx)
		var conditionJSON []byte = nil
		if proxy.Condition != nil {
			bytes, err := json.Marshal(proxy.Condition)
			if err != nil {
				return fmt.Errorf("failed to marshal condition: %w", err)
			}
			conditionJSON = bytes // Если есть данные, присваиваем []byte
		}
		err = repo.CreateProxy(ctx, &CreateProxyParams{
			ID:        uuid.New().String(),
			CreatedAt: pgtype.Timestamptz{Time: time.Now()},
			UpdatedAt: pgtype.Timestamptz{Time: time.Now()},
			ListenUrl: proxy.ListenURL,
			Mode:      string(proxy.Mode),
			PathKey:   proxy.PathKey,
			Tags:      proxy.Tags,
			Condition: conditionJSON,
		})

		if err != nil {
			return fmt.Errorf(
				"failed to insert proxy: %w, condition: %s, condition type: %T, condition value: %+v", err, conditionJSON, proxy.Condition, proxy.Condition,
			)
		}

		// Insert targets
		for i := range proxy.Targets {
			target := &proxy.Targets[i]

			err = repo.CreateTarget(ctx, &CreateTargetParams{
				ID:       target.ID,
				ProxyID:  target.ProxyID,
				Url:      target.URL,
				Weight:   target.Weight,
				IsActive: target.IsActive,
			})
			if err != nil {
				return fmt.Errorf("failed to insert target: %w", err)
			}
		}

		return nil
	})

	return err
}
