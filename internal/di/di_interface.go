package di

import userService "github.com/Haevnen/p2m_be/internal/pkg/user_service"

type Getter interface {
	RepositoryDI
	UseCaseDI
}

type RepositoryDI interface {
	UserRepository() userService.Repository
}

type UseCaseDI interface {
	UserUseCase() userService.UseCase
}
