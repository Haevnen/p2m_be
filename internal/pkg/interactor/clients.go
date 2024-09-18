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

	// validate client
	// if user has id and is_active is true, return error
	clientCheck, err := c.WithContext(ctx).Where(c.ClientID.Eq(clientDb.ClientID), c.IsActive).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if clientCheck != nil {
		return clientCheck.FromClient(), apperror.ErrClientHasIDExists
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

func (ci *ClientManagement) UpdateClient(ctx context.Context, clientID string, body p2mapi.UpdateClientBody) error {
	c := dal.Q.Client

	client, err := c.WithContext(ctx).Where(c.ClientID.Eq(clientID)).Where(c.IsActive.Is(true)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrRecordNotFound
		}
		return err
	}

	if body.EditingStyle != nil {
		client.EditingStyle = *body.EditingStyle
	}

	if body.Others != nil {
		client.Others = *body.Others
	}

	if body.Requirements != nil {
		client.Requirements = *body.Requirements
	}

	_, err = c.WithContext(ctx).Where(c.ClientID.Eq(clientID)).Where(c.IsActive.Is(true)).UpdateColumns(&client)
	return err
}

func (ci *ClientManagement) GetSingleClient(ctx context.Context, clientID string) (*p2mapi.ClientResponse, error) {
	c := dal.Q.Client

	client, err := c.WithContext(ctx).Where(c.ClientID.Eq(clientID)).Where(c.IsActive.Is(true)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}

	return client.FromClient(), nil
}
