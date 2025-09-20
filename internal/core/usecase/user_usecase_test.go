package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/core/domain"
	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/core/dto"
	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/core/usecase"
)

func (s *UserUsecaseSuiteTest) TestUserUseCase_Register() {
	tests := []struct {
		name        string
		input       dto.RegisterInput
		setupMocks  func()
		checkResult func(*testing.T, *dto.RegisterOutput, error)
	}{
		{
			name: "should register user successfully",
			input: dto.RegisterInput{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMocks: func() {
				// Check if email doesn't exist
				s.mockRepo.EXPECT().
					GetByEmail(s.ctx, "john@example.com").
					Return(nil, nil)
				// Create user
				s.mockRepo.EXPECT().
					Create(s.ctx, gomock.Any()).
					DoAndReturn(func(ctx interface{}, user *domain.User) error {
						user.UserID = 1 // Simulate repository assigning ID
						return nil
					})
			},
			checkResult: func(t *testing.T, output *dto.RegisterOutput, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, int64(1), output.UserID)
				assert.Equal(t, "John Doe", output.Name)
				assert.Equal(t, "john@example.com", output.Email)
			},
		},
		{
			name: "should return error when email already exists",
			input: dto.RegisterInput{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMocks: func() {
				s.mockRepo.EXPECT().
					GetByEmail(s.ctx, "john@example.com").
					Return(s.mockUsers[0], nil)
			},
			checkResult: func(t *testing.T, output *dto.RegisterOutput, err error) {
				assert.Error(t, err)
				assert.Nil(t, output)
				assert.Equal(t, usecase.ErrEmailAlreadyExists, err)
			},
		},
		{
			name: "should return error when input is invalid - empty name",
			input: dto.RegisterInput{
				Name:     "",
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMocks: func() {
				// No mock calls expected
			},
			checkResult: func(t *testing.T, output *dto.RegisterOutput, err error) {
				assert.Error(t, err)
				assert.Nil(t, output)
				assert.Equal(t, usecase.ErrInvalidInput, err)
			},
		},
		{
			name: "should return error when input is invalid - empty email",
			input: dto.RegisterInput{
				Name:     "John Doe",
				Email:    "",
				Password: "password123",
			},
			setupMocks: func() {
				// No mock calls expected
			},
			checkResult: func(t *testing.T, output *dto.RegisterOutput, err error) {
				assert.Error(t, err)
				assert.Nil(t, output)
				assert.Equal(t, usecase.ErrInvalidInput, err)
			},
		},
		{
			name: "should return error when input is invalid - empty password",
			input: dto.RegisterInput{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "",
			},
			setupMocks: func() {
				// No mock calls expected
			},
			checkResult: func(t *testing.T, output *dto.RegisterOutput, err error) {
				assert.Error(t, err)
				assert.Nil(t, output)
				assert.Equal(t, usecase.ErrInvalidInput, err)
			},
		},
		{
			name: "should return error when repository create fails",
			input: dto.RegisterInput{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMocks: func() {
				s.mockRepo.EXPECT().
					GetByEmail(s.ctx, "john@example.com").
					Return(nil, nil)
				s.mockRepo.EXPECT().
					Create(s.ctx, gomock.Any()).
					Return(assert.AnError)
			},
			checkResult: func(t *testing.T, output *dto.RegisterOutput, err error) {
				assert.Error(t, err)
				assert.Nil(t, output)
				assert.Equal(t, assert.AnError, err)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			output, err := s.useCase.Register(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, output, err)
		})
	}
}

func (s *UserUsecaseSuiteTest) TestUserUseCase_Login() {
	tests := []struct {
		name        string
		input       dto.LoginInput
		setupMocks  func()
		checkResult func(*testing.T, *dto.LoginOutput, error)
	}{
		{
			name: "should login successfully",
			input: dto.LoginInput{
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMocks: func() {
				// Mock user with hashed password for "password123"
				hashedPassword := "$2a$12$5CEGdJIUSFrHCyrSOPVEE.mdHjVucN38e2xRzCb8zM1XAB7ZfqdTS" // bcrypt hash for "password123"
				user := &domain.User{
					UserID:   1,
					Name:     "John Doe",
					Email:    "john@example.com",
					Password: hashedPassword,
				}
				s.mockRepo.EXPECT().
					GetByEmail(s.ctx, "john@example.com").
					Return(user, nil)
				s.mockJWTSigner.EXPECT().
					Sign(int64(1)).
					Return("jwt-token", nil)
			},
			checkResult: func(t *testing.T, output *dto.LoginOutput, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, "jwt-token", output.Token)
			},
		},
		{
			name: "should return error when input is invalid - empty email",
			input: dto.LoginInput{
				Email:    "",
				Password: "password123",
			},
			setupMocks: func() {
				// No mock calls expected
			},
			checkResult: func(t *testing.T, output *dto.LoginOutput, err error) {
				assert.Error(t, err)
				assert.Nil(t, output)
				assert.Equal(t, usecase.ErrInvalidInput, err)
			},
		},
		{
			name: "should return error when input is invalid - empty password",
			input: dto.LoginInput{
				Email:    "john@example.com",
				Password: "",
			},
			setupMocks: func() {
				// No mock calls expected
			},
			checkResult: func(t *testing.T, output *dto.LoginOutput, err error) {
				assert.Error(t, err)
				assert.Nil(t, output)
				assert.Equal(t, usecase.ErrInvalidInput, err)
			},
		},
		{
			name: "should return error when user not found",
			input: dto.LoginInput{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			setupMocks: func() {
				s.mockRepo.EXPECT().
					GetByEmail(s.ctx, "nonexistent@example.com").
					Return(nil, nil)
			},
			checkResult: func(t *testing.T, output *dto.LoginOutput, err error) {
				assert.Error(t, err)
				assert.Nil(t, output)
				assert.Equal(t, usecase.ErrInvalidCredentials, err)
			},
		},
		{
			name: "should return error when password is incorrect",
			input: dto.LoginInput{
				Email:    "john@example.com",
				Password: "wrongpassword",
			},
			setupMocks: func() {
				hashedPassword := "$2a$12$5CEGdJIUSFrHCyrSOPVEE.mdHjVucN38e2xRzCb8zM1XAB7ZfqdTS" // bcrypt hash for "password123"
				user := &domain.User{
					UserID:   1,
					Name:     "John Doe",
					Email:    "john@example.com",
					Password: hashedPassword,
				}
				s.mockRepo.EXPECT().
					GetByEmail(s.ctx, "john@example.com").
					Return(user, nil)
			},
			checkResult: func(t *testing.T, output *dto.LoginOutput, err error) {
				assert.Error(t, err)
				assert.Nil(t, output)
				assert.Equal(t, usecase.ErrInvalidCredentials, err)
			},
		},
		{
			name: "should return error when JWT signing fails",
			input: dto.LoginInput{
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMocks: func() {
				hashedPassword := "$2a$12$5CEGdJIUSFrHCyrSOPVEE.mdHjVucN38e2xRzCb8zM1XAB7ZfqdTS" // bcrypt hash for "password123"
				user := &domain.User{
					UserID:   1,
					Name:     "John Doe",
					Email:    "john@example.com",
					Password: hashedPassword,
				}
				s.mockRepo.EXPECT().
					GetByEmail(s.ctx, "john@example.com").
					Return(user, nil)
				s.mockJWTSigner.EXPECT().
					Sign(int64(1)).
					Return("", assert.AnError)
			},
			checkResult: func(t *testing.T, output *dto.LoginOutput, err error) {
				assert.Error(t, err)
				assert.Nil(t, output)
				assert.Equal(t, assert.AnError, err)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			output, err := s.useCase.Login(s.ctx, tt.input)

			// Assert
			tt.checkResult(t, output, err)
		})
	}
}

func (s *UserUsecaseSuiteTest) TestUserUseCase_GetMe() {
	tests := []struct {
		name        string
		userID      int64
		setupMocks  func()
		checkResult func(*testing.T, *dto.GetMeOutput, error)
	}{
		{
			name:   "should get user successfully",
			userID: 1,
			setupMocks: func() {
				s.mockRepo.EXPECT().
					GetByID(s.ctx, int64(1)).
					Return(s.mockUsers[0], nil)
			},
			checkResult: func(t *testing.T, output *dto.GetMeOutput, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, int64(1), output.UserID)
				assert.Equal(t, "John Doe", output.Name)
				assert.Equal(t, "john@example.com", output.Email)
			},
		},
		{
			name:   "should return error when userID is invalid",
			userID: 0,
			setupMocks: func() {
				// No mock calls expected
			},
			checkResult: func(t *testing.T, output *dto.GetMeOutput, err error) {
				assert.Error(t, err)
				assert.Nil(t, output)
				assert.Equal(t, usecase.ErrInvalidUserID, err)
			},
		},
		{
			name:   "should return error when user not found",
			userID: 999,
			setupMocks: func() {
				s.mockRepo.EXPECT().
					GetByID(s.ctx, int64(999)).
					Return(nil, nil)
			},
			checkResult: func(t *testing.T, output *dto.GetMeOutput, err error) {
				assert.Error(t, err)
				assert.Nil(t, output)
				assert.Equal(t, usecase.ErrUserNotFound, err)
			},
		},
		{
			name:   "should return error when repository fails",
			userID: 1,
			setupMocks: func() {
				s.mockRepo.EXPECT().
					GetByID(s.ctx, int64(1)).
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, output *dto.GetMeOutput, err error) {
				assert.Error(t, err)
				assert.Nil(t, output)
				assert.Equal(t, usecase.ErrUserNotFound, err)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			output, err := s.useCase.GetMe(s.ctx, tt.userID)

			// Assert
			tt.checkResult(t, output, err)
		})
	}
}