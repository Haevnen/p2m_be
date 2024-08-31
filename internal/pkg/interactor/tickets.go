package interactor

import (
	"context"
	"fmt"

	p2m_api "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/apperror"
	"github.com/Haevnen/p2m_be/internal/pkg/dal"
	"github.com/Haevnen/p2m_be/internal/pkg/model"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
	"github.com/go-openapi/swag"
)

type TicketManagement struct {
	userManagement   interactorinterface.UserManagementInterface
	clientManagement interactorinterface.ClientManagementInterface
	txManager        interactorinterface.TxManager
}

func NewTicketManagement(
	userManagement interactorinterface.UserManagementInterface,
	clientClientManagement interactorinterface.ClientManagementInterface,
	txManager interactorinterface.TxManager,
) *TicketManagement {
	return &TicketManagement{
		userManagement:   userManagement,
		clientManagement: clientClientManagement,
		txManager:        txManager,
	}
}

func (t *TicketManagement) AddTicketManual(ctx context.Context, body p2m_api.CreateTicketBody) error {
	payload := ctx.Value(model.AuthorizationPayloadKey).(*interactorinterface.Payload)
	if payload == nil {
		return apperror.ErrUserNotExists
	}

	// Get all users
	users, err := t.userManagement.GetAllUser(ctx, swag.Bool(true))
	if err != nil {
		return err
	}

	nickNameMapping := make(map[string]*p2m_api.User)

	for _, user := range users {
		nickNameMapping[user.NickName] = user
	}

	// Get client_id from client_name
	client, err := t.clientManagement.GetSingleClient(ctx, body.ClientId)
	if err != nil {
		return err
	}

	newTicket := model.Ticket{
		ClientID:    client.Id,
		Title:       body.Title,
		Description: body.Description,
		CreatedBy:   string(p2m_api.MANUAL),
		IsActive:    true,
		Status:      string(p2m_api.BACKLOG),
		QcID:        nickNameMapping[body.QcName].UserId,
		EditorID:    nickNameMapping[body.EditorName].UserId,
		Priority:    string(p2m_api.NORMAL),
	}

	// Start transaction to create new ticket and add record to history table
	return t.txManager.TransactionExec(ctx, func(childCtx context.Context) error {
		tx := childCtx.Value(txTransactionKey).(*dal.QueryTx)

		err := tx.Ticket.WithContext(childCtx).Create(&newTicket)
		if err != nil {
			return err
		}

		return tx.History.WithContext(childCtx).Create(&model.History{
			TicketID:    newTicket.ID,
			Action:      fmt.Sprintf("Ticket is created by %s", string(p2m_api.MANUAL)),
			PerformedBy: payload.UserID,
		})
	})
}

func (t *TicketManagement) UpdateTicket(ctx context.Context, ticketID int64, body p2m_api.UpdateTicketBody) error {
	return nil
}
