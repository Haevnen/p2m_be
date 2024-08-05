package usecase

import (
	"github.com/Haevnen/p2m_be/internal/pkg/base_service"
	"github.com/Haevnen/p2m_be/internal/pkg/base_service/usecase"
	"github.com/Haevnen/p2m_be/internal/pkg/model"
	userService "github.com/Haevnen/p2m_be/internal/pkg/user_service"
)

type UserServiceUseCase struct {
	*usecase.BaseServiceUseCase[model.User, int32]
	userRepository userService.Repository
}

func NewUserServiceUseCase(r userService.Repository) *UserServiceUseCase {
	return &UserServiceUseCase{
		BaseServiceUseCase: usecase.NewBaseServiceUseCase(interface{}(r).(base_service.Repository[model.User, int32])),
		userRepository:     r,
	}
}
