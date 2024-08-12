package interactor

import (
	"context"
	"errors"

	"gorm.io/gorm"

	p2mapi "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/apperror"
	"github.com/Haevnen/p2m_be/internal/pkg/dal"
	"github.com/Haevnen/p2m_be/internal/pkg/model"
)

type ClientManagement struct {
}

func NewClientManagement() *ClientManagement {
	return &ClientManagement{}
}

func (ci *ClientManagement) CreateClient(ctx context.Context, client p2mapi.ClientBody) (*p2mapi.ClientResponse, error) {
	c := dal.Q.Client

	var clientDb model.Client
	err := clientDb.ToClient(client)
	if err != nil {
		return nil, err
	}

	err = c.WithContext(ctx).Save(&clientDb)
	if err != nil {
		return nil, err
	}

	return clientDb.FromClient(), nil
}

func (ci *ClientManagement) RemoveClient(ctx context.Context, id string) error {
	c := dal.Q.Client

	client, err := c.WithContext(ctx).Where(c.ClientID.Eq(id), c.IsActive).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrRecordNotFound
		}

		return err
	}

	info, err := c.WithContext(ctx).Where(c.ClientID.Eq(client.ClientID)).Update(c.IsActive, false)
	if err != nil {
		return err
	}

	if info.RowsAffected != 1 {
		return apperror.ErrInternalServer
	}

	return nil
}

func (ci *ClientManagement) GetAllClient(ctx context.Context, includeDeActive *bool) ([]*p2mapi.ClientResponse, error) {
	c := dal.Q.Client
	var clients []*model.Client
	var err error

	if includeDeActive == nil || *includeDeActive == false {
		clients, err = c.WithContext(ctx).Where(c.IsActive).Find()
	} else {
		clients, err = c.WithContext(ctx).Find()
	}

	if err != nil {
		return nil, err
	}

	res := make([]*p2mapi.ClientResponse, 0, len(clients))
	for _, client := range clients {
		res = append(res, client.FromClient())
	}

	return res, nil
}
