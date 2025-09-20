package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/core/dto"
)

type UserController interface {
	Register(ctx context.Context, p Presenter, in dto.RegisterInput) ([]byte, error)
	Login(ctx context.Context, p Presenter, in dto.LoginInput) ([]byte, error)
	GetMe(ctx context.Context, p Presenter, userID int64) ([]byte, error)
}