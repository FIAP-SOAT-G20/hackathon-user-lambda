package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/core/dto"
)

type UserUseCase interface {
	Register(ctx context.Context, in dto.RegisterInput) (*dto.RegisterOutput, error)
	Login(ctx context.Context, in dto.LoginInput) (*dto.LoginOutput, error)
	GetMe(ctx context.Context, userID int64) (*dto.GetMeOutput, error)
}