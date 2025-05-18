package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/rakh1mbayev/Gym-Management-System/user_service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Repo
type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) CreateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepo) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepo) GetUserByToken(ctx context.Context, token string) (*domain.User, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepo) MarkEmailConfirmed(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockUserRepo) SetConfirmationToken(ctx context.Context, userID int64, token string) error {
	args := m.Called(ctx, userID, token)
	return args.Error(0)
}

func (m *mockUserRepo) ConfirmUserByToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

type mockMailService struct {
	mock.Mock
}

func (m *mockMailService) SendConfirmationEmail(to, subject, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

type mockNatsClient struct {
	mock.Mock
}

func (m *mockNatsClient) PublishUserRegistered(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func TestRegisterUser_Success(t *testing.T) {
	mockRepo := new(mockUserRepo)
	mockMailer := new(mockMailService)
	mockNats := new(mockNatsClient)

	uc := NewUserUsecase(mockRepo, mockMailer, mockNats)

	user := &domain.User{
		UserID:   1,
		Name:     "Alice",
		Email:    "alice@test.com",
		Password: "pass123",
		Role:     "user",
	}

	mockRepo.On("CreateUser", mock.Anything, user).Return(nil)
	mockRepo.On("SetConfirmationToken", mock.Anything, user.UserID, mock.Anything).Return(nil)
	mockNats.On("PublishUserRegistered", user.Email).Return(nil)
	mockMailer.On("SendConfirmationEmail", user.Email, mock.Anything, mock.Anything).Return(nil)

	err := uc.Register(context.Background(), user)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockMailer.AssertExpectations(t)
	mockNats.AssertExpectations(t)
}

func TestRegisterUser_FailToCreate(t *testing.T) {
	mockRepo := new(mockUserRepo)
	mockMailer := new(mockMailService)
	mockNats := new(mockNatsClient)

	uc := NewUserUsecase(mockRepo, mockMailer, mockNats)

	user := &domain.User{
		Name:     "Bob",
		Email:    "bob@test.com",
		Password: "pass123",
	}

	mockRepo.On("CreateUser", mock.Anything, user).Return(errors.New("db error"))

	err := uc.Register(context.Background(), user)
	assert.Error(t, err)

	mockMailer.AssertNotCalled(t, "SendConfirmationEmail", mock.Anything, mock.Anything, mock.Anything)
	mockNats.AssertNotCalled(t, "PublishUserRegistered", mock.Anything)
}

func TestSendEmail(t *testing.T) {
	mockRepo := new(mockUserRepo)
	mockMailer := new(mockMailService)
	mockNats := new(mockNatsClient)

	uc := NewUserUsecase(mockRepo, mockMailer, mockNats)
	user := &domain.User{
		UserID:   1,
		Name:     "Alice",
		Email:    "alice@test.com",
		Password: "secret",
		Role:     "user",
	}

	mockRepo.On("CreateUser", mock.Anything, user).Return(nil)
	mockRepo.On("SetConfirmationToken", mock.Anything, user.UserID, mock.Anything).Return(nil)
	mockNats.On("PublishUserRegistered", user.Email).Return(nil)
	mockMailer.On("SendConfirmationEmail", user.Email, mock.Anything, mock.Anything).Return(nil)

	err := uc.Register(context.Background(), user)
	assert.NoError(t, err)

	mockMailer.AssertCalled(t, "SendConfirmationEmail", user.Email, mock.Anything, mock.Anything)
}
