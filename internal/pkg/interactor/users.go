package interactor

import (
	"context"

	"github.com/Haevnen/p2m_be/internal/pkg/dal"
	"github.com/Haevnen/p2m_be/internal/pkg/model"
)

type UserManagement struct{}

func NewUserManagement() *UserManagement {
	return &UserManagement{}
}

func (ci *UserManagement) GetAllUser(ctx context.Context) (*model.User, error) {
	user, err := dal.Q.User.WithContext(ctx).First()
	if err != nil {
		return nil, err
	}

	return user, nil
}
