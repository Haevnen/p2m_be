package bank_service

import (
	"github.com/Haevnen/p2m_be/internal/pkg/base_service"
	"github.com/Haevnen/p2m_be/internal/pkg/model"
)

type UseCase interface {
	base_service.UseCase[model.User, int32]
}
