package registry

import (
	"github.com/Haevnen/p2m_be/internal/pkg/interactor"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
)

type Registry struct{}

func New() *Registry {
	return &Registry{}
}

func (r *Registry) UserManagementInteractor() interactorinterface.UserManagementInterface {
	return interactor.NewUserManagement()
}
