package repository

import (
	"github.com/Haevnen/p2m_be/internal/pkg/base_service/repository"
	"github.com/Haevnen/p2m_be/internal/pkg/db_wrap"
	"github.com/Haevnen/p2m_be/internal/pkg/model"
)

type MysqlUserServiceRepository struct {
	*repository.MysqlBaseServiceRepository[model.User, int32]
}

func NewMysqlUserServiceRepository(db db_wrap.DBGetter) *MysqlUserServiceRepository {
	return &MysqlUserServiceRepository{
		MysqlBaseServiceRepository: repository.NewMysqlBaseServiceRepository[model.User, int32](db),
	}
}
