package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ab-testing-service/internal/models"
	"github.com/ab-testing-service/internal/proxy"
)

type Storage struct {
	q     Querier
	db    *pgxpool.Pool
	Redis *redis.Client
}

const (
	proxyTTL  = 1 * time.Hour
	targetTTL = 1 * time.Hour
)

func NewStorage(conn *pgxpool.Pool, redis *redis.Client) *Storage {
	return &Storage{
		q:     New(conn),
		db:    conn,
		Redis: redis,
	}
}

func (s *Storage) SaveVisit(ctx context.Context, visit *models.Visit) error {
	visit.ID = uuid.New().String()
	visit.CreatedAt = time.Now()

	_, err := s.db.Exec(ctx,
		`INSERT INTO visits (id, proxy_id, target_id, user_id, rid, rrid, ruid, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		visit.ID, visit.ProxyID, visit.TargetID, visit.UserID,
		visit.RID, visit.RRID, visit.RUID, visit.CreatedAt,
	)
	return err
}

func (s *Storage) GetProxies(ctx context.Context) ([]proxy.Config, error) {
	var proxies []proxy.Config
	rows, err := s.db.Query(ctx,
		`SELECT id, listen_url, mode, condition, tags, path_key
		FROM proxies ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query proxies: %w", err)
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate proxies: %w", err)
	}

	for rows.Next() {
		var p models.Proxy
		var conditionJSON []byte
		if err := rows.Scan(&p.ID, &p.ListenURL, &p.Mode, &conditionJSON, &p.Tags, &p.PathKey); err != nil {
			return nil, fmt.Errorf("failed to scan proxy: %w", err)
		}
		if len(conditionJSON) > 0 {
			p.Condition = &models.RouteCondition{}
			if err := json.Unmarshal(conditionJSON, p.Condition); err != nil {
				return nil, fmt.Errorf("failed to unmarshal condition: %w", err)
			}
		}

		config := proxy.Config{
			ID:        p.ID,
			ListenURL: p.ListenURL,
			Mode:      p.Mode,
			Tags:      p.Tags,
		}

		if p.PathKey != nil {
			config.PathKey = *p.PathKey
		}

		condition, err := convertCondition(p.Condition)
		if err != nil {
			// Обработка ошибки
			log.Printf("Failed to convert condition for proxy %s: %v", p.ID, err)
			// Возможные варианты:
			// 1. Пропустить этот прокси
			//continue
			// 2. Вернуть ошибку выше
			//return nil, fmt.Errorf("failed to process proxy %s: %w", p.ID, err)
			// 3. Вернуть nil и продолжить обработку остальных прокси
			//config.Condition = nil
		}

		if condition != nil {
			config.Condition = condition
		}
		proxies = append(proxies, config)
	}
	return proxies, nil
}

// Безопасное приведение типов с обработкой ошибок
func convertCondition(rc *models.RouteCondition) (*proxy.Condition, error) {
	if rc == nil {
		return nil, nil // Если входной параметр nil, возвращаем nil без ошибки
	}

	// Проверяем поля на корректность
	if !rc.Type.IsValid() {
		return nil, fmt.Errorf("invalid condition type: %v", rc.Type)
	}

	// Создаем новый объект Condition
	condition := &proxy.Condition{
		Type:      rc.Type,
		ParamName: rc.ParamName,
		Values:    make(map[string]string),
		Default:   rc.Default,
	}

	// Копируем значения map, проверяя их валидность
	for k, v := range rc.Values {
		if k == "" {
			return nil, fmt.Errorf("empty key in Values map")
		}
		condition.Values[k] = v
	}

	return condition, nil
}

func convertConditionSafe(rc *models.RouteCondition) (c *proxy.Condition, err error) {
	// Защита от паники при конвертации
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to convert condition: %v", r)
		}
	}()

	if rc == nil {
		return nil, nil
	}

	return &proxy.Condition{
		Type:      rc.Type,
		ParamName: rc.ParamName,
		Values:    rc.Values,
		Default:   rc.Default,
	}, nil
}
