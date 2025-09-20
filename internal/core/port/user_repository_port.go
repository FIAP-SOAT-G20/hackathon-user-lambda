package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/core/domain"
)

type UserRepository interface {
	Create(ctx context.Context, u *domain.User) error
	GetByID(ctx context.Context, userID int64) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}