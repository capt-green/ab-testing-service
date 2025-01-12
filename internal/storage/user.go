package storage

import (
	"context"

	"github.com/ab-testing-service/internal/models"
)

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Storage) UserExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := s.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`,
		email,
	).Scan(&exists)
	return exists, err
}

func (s *Storage) CreateUser(ctx context.Context, user *models.User) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO users (id, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`,
		user.ID, user.Email, user.Password, user.CreatedAt, user.UpdatedAt,
	)
	return err
}
