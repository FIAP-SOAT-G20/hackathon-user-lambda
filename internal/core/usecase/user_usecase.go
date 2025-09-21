package usecase

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/core/domain"
	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/core/dto"
	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/core/port"
)

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidUserID      = errors.New("invalid user id")
	ErrUserNotFound       = errors.New("user not found")
)

type userUseCase struct {
	repo      port.UserRepository
	jwtSigner port.JWTSigner
}

func NewUserUseCase(repo port.UserRepository, jwtSigner port.JWTSigner) port.UserUseCase {
	return &userUseCase{repo: repo, jwtSigner: jwtSigner}
}

func (u *userUseCase) Register(ctx context.Context, in dto.RegisterInput) (*dto.RegisterOutput, error) {
	if in.Name == "" || in.Email == "" || in.Password == "" {
		return nil, ErrInvalidInput
	}

	// check if email already exists
	if existing, _ := u.repo.GetByEmail(ctx, in.Email); existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	user := &domain.User{
		// UserID will be assigned by repository (sequential)
		Name:      in.Name,
		Email:     in.Email,
		Password:  string(hash),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := u.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &dto.RegisterOutput{UserID: user.UserID, Name: user.Name, Email: user.Email}, nil
}

func (u *userUseCase) Login(ctx context.Context, in dto.LoginInput) (*dto.LoginOutput, error) {
	if in.Email == "" || in.Password == "" {
		return nil, ErrInvalidInput
	}
	user, err := u.repo.GetByEmail(ctx, in.Email)
	if err != nil || user == nil {
		return nil, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	token, err := u.jwtSigner.Sign(user.UserID)
	if err != nil {
		return nil, err
	}
	return &dto.LoginOutput{Token: token}, nil
}

func (u *userUseCase) GetMe(ctx context.Context, userID int64) (*dto.GetMeOutput, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}
	user, err := u.repo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}
	return &dto.GetMeOutput{UserID: user.UserID, Name: user.Name, Email: user.Email}, nil
}
