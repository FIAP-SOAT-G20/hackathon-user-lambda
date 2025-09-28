package controller

import (
	"context"

	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/core/dto"
	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/core/port"
)

type UserController struct {
	usecase port.UserUseCase
}

func NewUserController(uc port.UserUseCase) port.UserController {
	return &UserController{usecase: uc}
}

func (c *UserController) Register(ctx context.Context, p port.Presenter, in dto.RegisterInput) ([]byte, error) {
	out, err := c.usecase.Register(ctx, in)
	if err != nil {
		return nil, err
	}
	return p.Present(out)
}

func (c *UserController) Login(ctx context.Context, p port.Presenter, in dto.LoginInput) ([]byte, error) {
	out, err := c.usecase.Login(ctx, in)
	if err != nil {
		return nil, err
	}
	return p.Present(out)
}

func (c *UserController) GetMe(ctx context.Context, p port.Presenter, userID int64) ([]byte, error) {
	out, err := c.usecase.GetMe(ctx, userID)
	if err != nil {
		return nil, err
	}
	return p.Present(out)
}

func (c *UserController) GetUserByID(ctx context.Context, p port.Presenter, userID int64) ([]byte, error) {
	out, err := c.usecase.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return p.Present(out)
}
