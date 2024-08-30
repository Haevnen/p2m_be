package interactor

import (
	"context"

	p2mapi "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/apperror"
	"github.com/Haevnen/p2m_be/internal/pkg/model"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
	"github.com/go-openapi/swag"
)

type TicketManagement struct {
	userManagement   interactorinterface.UserManagementInterface
	clientManagement interactorinterface.ClientManagementInterface
}

func NewTicketManagement(
	userManagement interactorinterface.UserManagementInterface,
	clientClientManagement interactorinterface.ClientManagementInterface,
) *TicketManagement {
	return &TicketManagement{
		userManagement:   userManagement,
		clientManagement: clientClientManagement,
	}
}

func (t *TicketManagement) AddTicketManual(ctx context.Context, body p2mapi.CreateTicketBody) error {
	payload := ctx.Value(model.AuthorizationPayloadKey).(*interactorinterface.Payload)
	if payload == nil {
		return apperror.ErrUserNotExists
	}

	// Get all users
	users, err := t.userManagement.GetAllUser(ctx, swag.Bool(true))
	if err != nil {
		return err
	}

	nickNameMapping := make(map[string]*p2mapi.User)

	for _, user := range users {
		nickNameMapping[user.NickName] = user
	}

	// Get client_id from client_name
	client, err := t.clientManagement.GetSingleClient(ctx, body.ClientId)
	if err != nil {
		return err
	}

	newTicket := model.Ticket{
		ClientID: client.Id,
		Subject:  body.Subject,
		Content:  body.Content,
		UserID:   payload.UserID,
		NickName: payload.NickName,
	}

}
func (t *TicketManagement) UpdateTicket(ctx context.Context, ticketID int64, body p2mapi.UpdateTicketBody) error {

}
