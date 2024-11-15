package mock

import (
	"context"
	"errors"
	"maqhaa/product_service/external/model"
)

type MockUserRepository struct {
	GetUserFunc func(ctx context.Context, token string) (*model.UserData, error)

	userResponses map[string]*model.UserData
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		userResponses: make(map[string]*model.UserData),
	}
}

func (m *MockUserRepository) SetUserResponse(token string, user *model.UserData) {
	m.userResponses[token] = user
}

func (m *MockUserRepository) GetUser(ctx context.Context, token string) (*model.UserData, error) {

	if m.GetUserFunc != nil {
		return m.GetUserFunc(ctx, token)
	}

	if user, ok := m.userResponses[token]; ok {
		return user, nil
	}

	return nil, errors.New("GetUserFunc not implemented in the mock")
}
