package di

import (
	"gorm.io/gorm"

	"github.com/Haevnen/p2m_be/internal/pkg/db_wrap"
	userService "github.com/Haevnen/p2m_be/internal/pkg/user_service"
	userRepository "github.com/Haevnen/p2m_be/internal/pkg/user_service/repository"
	userUseCase "github.com/Haevnen/p2m_be/internal/pkg/user_service/usecase"
)

type DI struct {
	db db_wrap.DBGetter
}

func New(db *gorm.DB) *DI {
	return &DI{
		db: db_wrap.NewDBGetter(db),
	}
}

func (r *DI) UserRepository() userService.Repository {
	return userRepository.NewMysqlUserServiceRepository(r.db)
}

func (r *DI) UserUseCase() userService.UseCase {
	return userUseCase.NewUserServiceUseCase(r.UserRepository())
}
