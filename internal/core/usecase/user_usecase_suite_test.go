package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/core/domain"
	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/core/port"
	mockport "github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/core/port/mocks"
	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/core/usecase"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type UserUsecaseSuiteTest struct {
	suite.Suite
	mockUsers     []*domain.User
	mockRepo      *mockport.MockUserRepository
	mockJWTSigner *mockport.MockJWTSigner
	useCase       port.UserUseCase
	ctx           context.Context
	ctrl          *gomock.Controller
}

func (s *UserUsecaseSuiteTest) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockRepo = mockport.NewMockUserRepository(s.ctrl)
	s.mockJWTSigner = mockport.NewMockJWTSigner(s.ctrl)
	s.useCase = usecase.NewUserUseCase(s.mockRepo, s.mockJWTSigner)
	s.ctx = context.Background()
	currentTime := time.Now().Unix()
	s.mockUsers = []*domain.User{
		{
			UserID:    1,
			Name:      "John Doe",
			Email:     "john@example.com",
			Password:  "$2a$10$hashedpassword1",
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		},
		{
			UserID:    2,
			Name:      "Jane Smith",
			Email:     "jane@example.com",
			Password:  "$2a$10$hashedpassword2",
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		},
	}
}

func (s *UserUsecaseSuiteTest) TearDownTest() {
	s.ctrl.Finish()
}

func TestUserUsecaseSuiteTest(t *testing.T) {
	suite.Run(t, new(UserUsecaseSuiteTest))
}