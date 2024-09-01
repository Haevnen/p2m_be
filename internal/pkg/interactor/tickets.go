package interactor

import (
	"context"
	"errors"
	"fmt"

	p2m_api "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/apperror"
	"github.com/Haevnen/p2m_be/internal/pkg/dal"
	"github.com/Haevnen/p2m_be/internal/pkg/model"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
	"github.com/go-openapi/swag"
	"gorm.io/gorm"
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

	if _, ok := nickNameMapping[body.QcName]; !ok {
		return apperror.ErrUserNotExists
	}

	if _, ok := nickNameMapping[body.EditorName]; !ok {
		return apperror.ErrUserNotExists
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

		// Create new ticket
		err := tx.Ticket.WithContext(childCtx).Create(&newTicket)
		if err != nil {
			return err
		}

		// Create new link
		if body.Links != nil {
			var links []*model.Link
			for _, link := range *body.Links {
				links = append(links, &model.Link{
					TicketID: newTicket.ID,
					Link:     link,
				})
			}

			if len(links) > 0 {
				err = tx.Link.WithContext(childCtx).CreateInBatches(links, 50)
				if err != nil {
					return err
				}
			}
		}

		// Create new history
		return tx.History.WithContext(childCtx).Create(&model.History{
			TicketID:    newTicket.ID,
			Action:      fmt.Sprintf("Ticket is created by %s", string(p2m_api.MANUAL)),
			PerformedBy: payload.UserID,
		})
	})
}

func (t *TicketManagement) UpdateTicket(ctx context.Context, ticketID int64, body p2m_api.UpdateTicketBody) error {
	// Get payload to check if user is admin
	payload := ctx.Value(model.AuthorizationPayloadKey).(*interactorinterface.Payload)
	if payload == nil {
		return apperror.ErrUserNotExists
	}

	// Only admin can update Title, Description, ClientId and Done Status
	if body.Title != nil || body.Description != nil || body.ClientId != nil || (body.Status != nil && string(*(body.Status)) == string(p2m_api.DONE)) {
		if !payload.IsAdmin {
			return apperror.ErrUserNotExists
		}
	}

	u := dal.Q.Ticket

	ticket, err := u.WithContext(ctx).Where(u.ID.Eq(ticketID)).Where(u.IsActive.Is(true)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrTicketNotFound
		}
		return err
	}

	// Get all users
	users, err := t.userManagement.GetAllUser(ctx, swag.Bool(true))
	if err != nil {
		return err
	}

	nickNameMapping := make(map[string]*p2m_api.User)
	userIdMapping := make(map[string]*p2m_api.User)

	for _, user := range users {
		nickNameMapping[user.NickName] = user
		userIdMapping[user.UserId] = user
	}

	var histories []*model.History

	if body.Title != nil {
		ticket.Title = *body.Title
		histories = append(histories, &model.History{
			TicketID: ticket.ID,
			Action:   fmt.Sprintf("User %s update ticket title", userIdMapping[payload.UserID].NickName),
		})
	}

	if body.Description != nil {
		ticket.Description = *body.Description
		histories = append(histories, &model.History{
			TicketID: ticket.ID,
			Action:   fmt.Sprintf("User %s update ticket description. ", userIdMapping[payload.UserID].NickName),
		})
	}

	if body.ClientId != nil {
		client, err := t.clientManagement.GetSingleClient(ctx, *body.ClientId)
		if err != nil {
			return err
		}
		ticket.ClientID = client.Id
		histories = append(histories, &model.History{
			TicketID: ticket.ID,
			Action:   fmt.Sprintf("User %s update ticket client. ", userIdMapping[payload.UserID].NickName),
		})
	}

	if body.Status != nil {
		ticket.Status = string(*(body.Status))
		histories = append(histories, &model.History{
			TicketID: ticket.ID,
			Action:   fmt.Sprintf("User %s update ticket status from %s to %s.", userIdMapping[payload.UserID].NickName, string(ticket.Status), string(*(body.Status))),
		})
	}

	if body.Priority != nil {
		ticket.Priority = string(*(body.Priority))
		histories = append(histories, &model.History{
			TicketID: ticket.ID,
			Action:   fmt.Sprintf("User %s update ticket priority from %s to %s.", userIdMapping[payload.UserID].NickName, string(ticket.Priority), string(*(body.Priority))),
		})
	}

	if body.QcName != nil {
		if _, ok := nickNameMapping[*body.QcName]; !ok {
			return apperror.ErrUserNotExists
		}

		ticket.QcID = nickNameMapping[*body.QcName].UserId
		histories = append(histories, &model.History{
			TicketID: ticket.ID,
			Action:   fmt.Sprintf("User %s update QC assignee from %s to %s.", userIdMapping[payload.UserID].NickName, userIdMapping[ticket.QcID].NickName, *body.QcName),
		})
	}

	if body.EditorName != nil {
		if _, ok := nickNameMapping[*body.EditorName]; !ok {
			return apperror.ErrUserNotExists
		}

		ticket.EditorID = nickNameMapping[*body.EditorName].UserId
		histories = append(histories, &model.History{
			TicketID: ticket.ID,
			Action:   fmt.Sprintf("User %s update Editor assignee from %s to %s.", userIdMapping[payload.UserID].NickName, userIdMapping[ticket.EditorID].NickName, *body.EditorName),
		})
	}

	return t.txManager.TransactionExec(ctx, func(childCtx context.Context) error {
		tx := childCtx.Value(txTransactionKey).(*dal.QueryTx)

		// Create histories
		if len(histories) > 0 {
			err = tx.History.WithContext(childCtx).CreateInBatches(histories, 50)
			if err != nil {
				return err
			}
		}

		_, err = tx.Ticket.WithContext(childCtx).Where(tx.Ticket.ID.Eq(ticketID)).Where(tx.Ticket.IsActive.Is(true)).UpdateColumns(&ticket)
		return err
	})
}
