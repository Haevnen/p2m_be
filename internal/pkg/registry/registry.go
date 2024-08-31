package registry

import (
	"github.com/Haevnen/p2m_be/internal/pkg/interactor"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
)

type Registry struct {
	key string
}

func New(key string) *Registry {
	return &Registry{
		key: key,
	}
}

func (r *Registry) UserManagementInteractor() interactorinterface.UserManagementInterface {
	return interactor.NewUserManagement(r.PasetoMaker())
}

func (r *Registry) ClientManagementInteractor() interactorinterface.ClientManagementInterface {
	return interactor.NewClientManagement()
}

func (r *Registry) LinkManagementInterface() interactorinterface.LinkManagementInterface {
	return interactor.NewLinkManagement()
}

func (r *Registry) CommentManagementInterface() interactorinterface.CommentManagementInterface {
	return interactor.NewCommentManagement()
}

func (r *Registry) HistoryManagementInterface() interactorinterface.HistoryManagementInterface {
	return interactor.NewHistoryManagement()
}

func (r *Registry) TicketManagementInterface() interactorinterface.TicketManagementInterface {
	return interactor.NewTicketManagement(r.UserManagementInteractor(), r.ClientManagementInteractor(), r.TxManager())
}

func (r *Registry) PasetoMaker() interactorinterface.Maker {
	paseto, err := interactor.NewPasetoMaker(r.key)
	if err != nil {
		panic(err)
	}
	return paseto
}

func (r *Registry) TxManager() interactorinterface.TxManager {
	return interactor.NewTxManager()
}
