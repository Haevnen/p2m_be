package interactorinterface

import (
	"context"

	"github.com/Haevnen/p2m_be/internal/pkg/model"
)

type UserManagementInterface interface {
	GetAllUser(context.Context) (*model.User, error)
}
