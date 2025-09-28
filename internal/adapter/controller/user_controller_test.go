package controller_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/adapter/controller"
	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/core/dto"
	mockport "github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/core/port/mocks"
)

func TestUserController_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mockport.NewMockUserUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	c := controller.NewUserController(mockUC)

	ctx := context.Background()
	in := dto.RegisterInput{Name: "Alice", Email: "a@a.com", Password: "123"}
	out := &dto.RegisterOutput{UserID: 1, Name: "Alice", Email: "a@a.com"}

	mockUC.EXPECT().Register(ctx, in).Return(out, nil)
	mockPresenter.EXPECT().Present(gomock.AssignableToTypeOf(&dto.RegisterOutput{})).Return([]byte("{}"), nil)

	b, err := c.Register(ctx, mockPresenter, in)
	assert.NoError(t, err)
	assert.NotNil(t, b)
}

func TestUserController_Register_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mockport.NewMockUserUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	c := controller.NewUserController(mockUC)

	ctx := context.Background()
	in := dto.RegisterInput{Name: "Alice", Email: "a@a.com", Password: "123"}

	mockUC.EXPECT().Register(ctx, in).Return(nil, assert.AnError)

	b, err := c.Register(ctx, mockPresenter, in)
	assert.Error(t, err)
	assert.Nil(t, b)
}

func TestUserController_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mockport.NewMockUserUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	c := controller.NewUserController(mockUC)

	ctx := context.Background()
	in := dto.LoginInput{Email: "a@a.com", Password: "123"}
	out := &dto.LoginOutput{Token: "abc"}

	mockUC.EXPECT().Login(ctx, in).Return(out, nil)
	mockPresenter.EXPECT().Present(gomock.AssignableToTypeOf(&dto.LoginOutput{})).Return([]byte("{}"), nil)

	b, err := c.Login(ctx, mockPresenter, in)
	assert.NoError(t, err)
	assert.NotNil(t, b)
}

func TestUserController_Login_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mockport.NewMockUserUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	c := controller.NewUserController(mockUC)

	ctx := context.Background()
	in := dto.LoginInput{Email: "a@a.com", Password: "123"}

	mockUC.EXPECT().Login(ctx, in).Return(nil, assert.AnError)

	b, err := c.Login(ctx, mockPresenter, in)
	assert.Error(t, err)
	assert.Nil(t, b)
}

func TestUserController_GetMe_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mockport.NewMockUserUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	c := controller.NewUserController(mockUC)

	ctx := context.Background()
	userID := int64(5)
	out := &dto.GetMeOutput{UserID: userID, Name: "Alice", Email: "a@a.com"}

	mockUC.EXPECT().GetMe(ctx, userID).Return(out, nil)
	mockPresenter.EXPECT().Present(gomock.AssignableToTypeOf(&dto.GetMeOutput{})).Return([]byte("{}"), nil)

	b, err := c.GetMe(ctx, mockPresenter, userID)
	assert.NoError(t, err)
	assert.NotNil(t, b)
}

func TestUserController_GetMe_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mockport.NewMockUserUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	c := controller.NewUserController(mockUC)

	ctx := context.Background()
	userID := int64(5)

	mockUC.EXPECT().GetMe(ctx, userID).Return(nil, assert.AnError)

	b, err := c.GetMe(ctx, mockPresenter, userID)
	assert.Error(t, err)
	assert.Nil(t, b)
}

func TestUserController_GetUserByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mockport.NewMockUserUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	c := controller.NewUserController(mockUC)

	ctx := context.Background()
	userID := int64(5)
	out := &dto.GetUserByIDOutput{UserID: userID, Name: "Alice", Email: "a@a.com"}

	mockUC.EXPECT().GetUserByID(ctx, userID).Return(out, nil)
	mockPresenter.EXPECT().Present(gomock.AssignableToTypeOf(&dto.GetUserByIDOutput{})).Return([]byte("{}"), nil)

	b, err := c.GetUserByID(ctx, mockPresenter, userID)
	assert.NoError(t, err)
	assert.NotNil(t, b)
}

func TestUserController_GetUserByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mockport.NewMockUserUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	c := controller.NewUserController(mockUC)

	ctx := context.Background()
	userID := int64(5)

	mockUC.EXPECT().GetUserByID(ctx, userID).Return(nil, assert.AnError)

	b, err := c.GetUserByID(ctx, mockPresenter, userID)
	assert.Error(t, err)
	assert.Nil(t, b)
}
