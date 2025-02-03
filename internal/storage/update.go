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

// UpdateProxyWithTargetsAndCondition обновляет прокси, его таргеты и условия в одной транзакции
func (s *Storage) UpdateProxyWithTargetsAndCondition(ctx context.Context, proxyID string,
	currentProxy *models.Proxy, targets []models.Target,
	condition *models.RouteCondition, userID *string) error {

	err := pgx.BeginFunc(ctx, s.db, func(tx pgx.Tx) (err error) {
		repo := New(tx)

		previousJSON, err := json.Marshal(currentProxy.Targets)
		if err != nil {
			return fmt.Errorf("failed to marshal previous state: %w", err)
		}

		newJSON, err := json.Marshal(targets)
		if err != nil {
			return fmt.Errorf("failed to marshal new state: %w", err)
		}

		if err = repo.CreateProxyChange(ctx, &CreateProxyChangeParams{
			ID:            uuid.New().String(),
			ProxyID:       proxyID,
			ChangeType:    string(models.ChangeTypeTargetsUpdate),
			PreviousState: previousJSON,
			NewState:      newJSON,
			CreatedAt:     pgtype.Timestamptz{Time: time.Now()},
			CreatedBy:     userID,
		}); err != nil {
			return fmt.Errorf("failed to record target changes: %w", err)
		}

		if condition == nil {
			return nil
		}

		previousJSON, err = json.Marshal(currentProxy.Condition)
		if err != nil {
			return fmt.Errorf("failed to marshal previous state: %w", err)
		}

		newJSON, err = json.Marshal(condition)
		if err != nil {
			return fmt.Errorf("failed to marshal new state: %w", err)
		}

		if err = repo.CreateProxyChange(ctx, &CreateProxyChangeParams{
			ID:            uuid.New().String(),
			ProxyID:       proxyID,
			ChangeType:    string(models.ChangeTypeConditionUpdate),
			PreviousState: previousJSON,
			NewState:      newJSON,
			CreatedAt:     pgtype.Timestamptz{Time: time.Now()},
			CreatedBy:     userID,
		}); err != nil {
			return fmt.Errorf("failed to record condition changes: %w", err)
		}

		if err = repo.DeleteTargetByProxyID(ctx, proxyID); err != nil {
			return fmt.Errorf("failed to delete targets: %w", err)
		}

		for _, target := range targets {
			if err = repo.CreateTarget(ctx, &CreateTargetParams{
				ID:       target.ID,
				Url:      target.URL,
				ProxyID:  proxyID,
				Weight:   target.Weight,
				IsActive: target.IsActive,
			}); err != nil {
				return fmt.Errorf("failed to insert target: %w", err)
			}
		}

		if condition == nil {
			return nil
		}

		conditionJSON, err := json.Marshal(condition)
		if err != nil {
			return fmt.Errorf("failed to marshal condition: %w", err)
		}

		if err = repo.UpdateProxyCondition(ctx, &UpdateProxyConditionParams{
			ID:        proxyID,
			Condition: conditionJSON,
		}); err != nil {
			return fmt.Errorf("failed to update condition: %w", err)
		}

		return nil
	})

	return err
}
