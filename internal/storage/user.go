package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/ab-testing-service/internal/models"
)

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	userModel := models.User{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.PasswordHash,
	}
	return &userModel, nil
}

func (s *Storage) UserExists(ctx context.Context, email string) (bool, error) {
	exists, err := s.q.UserExists(ctx, email)
	return exists, err
}

func (s *Storage) CreateUser(ctx context.Context, user *models.User) error {
	err := s.q.CreateUser(ctx, &CreateUserParams{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.Password,
		CreatedAt:    pgtype.Timestamptz{Time: user.CreatedAt},
		UpdatedAt:    pgtype.Timestamptz{Time: user.UpdatedAt},
	})
	return err
}
