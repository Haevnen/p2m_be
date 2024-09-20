package interactor

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-openapi/swag"
	"gorm.io/gorm"

	p2mapi "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/apperror"
	"github.com/Haevnen/p2m_be/internal/pkg/dal"
	"github.com/Haevnen/p2m_be/internal/pkg/model"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
	"github.com/Haevnen/p2m_be/pkg/constants"
	"github.com/Haevnen/p2m_be/pkg/util"
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

	if _, ok := nickNameMapping[body.QcName]; !ok {
		return apperror.ErrQCNameNotExists
	}

	if _, ok := nickNameMapping[body.EditorName]; !ok {
		return apperror.ErrEditorNameNotExists
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
		CreatedBy:   string(p2mapi.MANUAL),
		IsActive:    true,
		Status:      string(p2mapi.BACKLOG),
		QcID:        nickNameMapping[body.QcName].UserId,
		EditorID:    nickNameMapping[body.EditorName].UserId,
		Priority:    string(body.Priority),
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
			Action:      fmt.Sprintf("Ticket is created by %s", string(p2mapi.MANUAL)),
			PerformedBy: payload.UserID,
		})
	})
}

func (t *TicketManagement) UpdateTicket(ctx context.Context, ticketID int64, body p2mapi.UpdateTicketBody) error {
	// Get payload to check if user is admin
	payload := ctx.Value(model.AuthorizationPayloadKey).(*interactorinterface.Payload)
	if payload == nil {
		return apperror.ErrUserNotExists
	}

	// Only admin can update Title, Description, ClientId and Done Status
	if body.Title != nil || body.Description != nil || body.ClientId != nil || (body.Status != nil && string(*(body.Status)) == string(p2mapi.DONE)) {
		if !payload.IsAdmin {
			return apperror.ErrForbidden
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

	// Only admin user can update status from DONE to another
	if body.Status != nil && ticket.Status == string(p2mapi.DONE) {
		if !payload.IsAdmin {
			return apperror.ErrForbidden
		}
	}

	// Get all users
	users, err := t.userManagement.GetAllUser(ctx, swag.Bool(true))
	if err != nil {
		return err
	}

	nickNameMapping := make(map[string]*p2mapi.User)
	userIdMapping := make(map[string]*p2mapi.User)

	for _, user := range users {
		nickNameMapping[user.NickName] = user
		userIdMapping[user.UserId] = user
	}

	var histories []*model.History

	if body.Title != nil {
		ticket.Title = *body.Title
		histories = append(histories, &model.History{
			TicketID:    ticket.ID,
			Action:      fmt.Sprintf("User %s update ticket title", userIdMapping[payload.UserID].NickName),
			PerformedBy: payload.UserID,
		})
	}

	if body.Description != nil {
		ticket.Description = *body.Description
		histories = append(histories, &model.History{
			TicketID:    ticket.ID,
			Action:      fmt.Sprintf("User %s update ticket description. ", userIdMapping[payload.UserID].NickName),
			PerformedBy: payload.UserID,
		})
	}

	if body.ClientId != nil {
		client, err := t.clientManagement.GetSingleClient(ctx, *body.ClientId)
		if err != nil {
			return err
		}
		ticket.ClientID = client.Id
		histories = append(histories, &model.History{
			TicketID:    ticket.ID,
			Action:      fmt.Sprintf("User %s update ticket client. ", userIdMapping[payload.UserID].NickName),
			PerformedBy: payload.UserID,
		})
	}

	if body.Status != nil {
		histories = append(histories, &model.History{
			TicketID:    ticket.ID,
			Action:      fmt.Sprintf("User %s update ticket status from %s to %s.", userIdMapping[payload.UserID].NickName, ticket.Status, string(*(body.Status))),
			PerformedBy: payload.UserID,
		})
		ticket.Status = string(*(body.Status))
	}

	if body.NumOfSingleImage != nil {
		histories = append(histories, &model.History{
			TicketID:    ticket.ID,
			Action:      fmt.Sprintf("User %s update quantity of single image from %d to %d.", userIdMapping[payload.UserID].NickName, ticket.NumOfSingleImage, *body.NumOfSingleImage),
			PerformedBy: payload.UserID,
		})
		ticket.NumOfSingleImage = *body.NumOfSingleImage
	}

	if body.NumOfMultipleImage != nil {
		histories = append(histories, &model.History{
			TicketID:    ticket.ID,
			Action:      fmt.Sprintf("User %s update quantity of multiple image from %d to %d.", userIdMapping[payload.UserID].NickName, ticket.NumOfMultipleImage, *body.NumOfMultipleImage),
			PerformedBy: payload.UserID,
		})
		ticket.NumOfMultipleImage = *body.NumOfMultipleImage
	}

	if body.Priority != nil {
		histories = append(histories, &model.History{
			TicketID:    ticket.ID,
			Action:      fmt.Sprintf("User %s update ticket priority from %s to %s.", userIdMapping[payload.UserID].NickName, ticket.Priority, string(*(body.Priority))),
			PerformedBy: payload.UserID,
		})
		ticket.Priority = string(*body.Priority)
	}

	if body.QcName != nil {
		if _, ok := nickNameMapping[*body.QcName]; !ok {
			return apperror.ErrQCNameNotExists
		}

		histories = append(histories, &model.History{
			TicketID:    ticket.ID,
			Action:      fmt.Sprintf("User %s update QC assignee from %s to %s.", userIdMapping[payload.UserID].NickName, userIdMapping[ticket.QcID].NickName, *body.QcName),
			PerformedBy: payload.UserID,
		})
		ticket.QcID = nickNameMapping[*body.QcName].UserId
	}

	if body.EditorName != nil {
		if _, ok := nickNameMapping[*body.EditorName]; !ok {
			return apperror.ErrEditorNameNotExists
		}

		histories = append(histories, &model.History{
			TicketID:    ticket.ID,
			Action:      fmt.Sprintf("User %s update Editor assignee from %s to %s.", userIdMapping[payload.UserID].NickName, userIdMapping[ticket.EditorID].NickName, *body.EditorName),
			PerformedBy: payload.UserID,
		})
		ticket.EditorID = nickNameMapping[*body.EditorName].UserId
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

		_, err = tx.Ticket.WithContext(childCtx).Where(tx.Ticket.ID.Eq(ticketID)).Where(tx.Ticket.IsActive.Is(true)).Omit(tx.Ticket.UpdatedAt).Updates(&ticket)
		return err
	})
}

func (t *TicketManagement) GetAllTicketsByContractType(ctx context.Context) ([]*p2mapi.ListTicketItem, error) {
	ti := dal.Q.Ticket
	tid := ti.WithContext(ctx)
	u := dal.Q.User

	// get user to get contract type
	payload := ctx.Value(model.AuthorizationPayloadKey).(*interactorinterface.Payload)
	if payload == nil {
		return nil, apperror.ErrUserNotExists
	}

	user, err := u.WithContext(ctx).Where(u.UserID.Eq(payload.UserID)).Where(u.IsActive.Is(true)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrRecordNotFound
		}

		return nil, err
	}

	if user == nil {
		return nil, apperror.ErrUserNotExists
	}

	var ticketsDb []*model.Ticket
	now := time.Now()
	// general query
	// Only return needed fields (title, ticket-number, status, priority)
	ticketQuery := tid.Select(ti.Title, ti.ID, ti.Status, ti.Priority, ti.QcID, ti.EditorID, ti.UpdatedAt).
		// Ignore deleted ticket
		Where(ti.IsActive.Is(true)).
		// List all ticket not completed (for all day)
		Where(tid.Where(ti.Status.Neq(string(p2mapi.DONE))).
			// Or List all ticket in status DONE (for current day)
			Or(tid.Where(ti.UpdatedAt.Between(util.Begin(now), util.End(now))).Where(ti.Status.Eq(string(p2mapi.DONE)))))
	// List by priority and updated_at asc
	ticketQuery.Order(ti.Priority.Desc()).Order(ti.UpdatedAt.Desc())

	// if contract type is FREELANCE, return only ticket from themselves
	if user.ContractType == string(p2mapi.FREELANCE) {
		ticketsDb, err = ticketQuery.Where(tid.Where(ti.EditorID.Eq(user.UserID)).Or(ti.QcID.Eq(user.UserID))).Find()
		if err != nil {
			return nil, err
		}
	} else {
		ticketsDb, err = ticketQuery.Find()
		if err != nil {
			return nil, err
		}
	}

	tickets := []*p2mapi.ListTicketItem{{Status: p2mapi.BACKLOG, Tickets: make([]p2mapi.ListTicket, 0)},
		{Status: p2mapi.INPROGRESS, Tickets: make([]p2mapi.ListTicket, 0)},
		{Status: p2mapi.READYTOQC, Tickets: make([]p2mapi.ListTicket, 0)},
		{Status: p2mapi.QCVERIFYING, Tickets: make([]p2mapi.ListTicket, 0)},
		{Status: p2mapi.QCDONE, Tickets: make([]p2mapi.ListTicket, 0)},
		{Status: p2mapi.DONE, Tickets: make([]p2mapi.ListTicket, 0)}}

	// convert to api model
	if len(ticketsDb) > 0 {
		// Get all users
		users, err := t.userManagement.GetAllUser(ctx, swag.Bool(true))
		if err != nil {
			return nil, err
		}
		userIdMapping := make(map[string]*p2mapi.User)
		for _, userItem := range users {
			userIdMapping[userItem.UserId] = userItem
		}

		for _, ticket := range ticketsDb {
			ticketItem := p2mapi.ListTicket{
				EditorName: userIdMapping[ticket.EditorID].NickName,
				Id:         ticket.ID,
				Priority:   p2mapi.Priority(ticket.Priority),
				QcName:     userIdMapping[ticket.QcID].NickName,
				Title:      ticket.Title,
				UpdatedAt:  ticket.UpdatedAt.Format(constants.DateTimeFormat),
			}

			switch ticket.Status {
			case string(p2mapi.BACKLOG):
				tickets[0].Tickets = append(tickets[0].Tickets, ticketItem)
			case string(p2mapi.INPROGRESS):
				tickets[1].Tickets = append(tickets[1].Tickets, ticketItem)
			case string(p2mapi.READYTOQC):
				tickets[2].Tickets = append(tickets[2].Tickets, ticketItem)
			case string(p2mapi.QCVERIFYING):
				tickets[3].Tickets = append(tickets[3].Tickets, ticketItem)
			case string(p2mapi.QCDONE):
				tickets[4].Tickets = append(tickets[4].Tickets, ticketItem)
			case string(p2mapi.DONE):
				tickets[5].Tickets = append(tickets[5].Tickets, ticketItem)
			}
		}
	}

	return tickets, nil
}

func (t *TicketManagement) GetTicketById(ctx context.Context, ticketId int64) (*p2mapi.SingleTicketResponse, error) {
	ti := dal.Q.Ticket
	tid := ti.WithContext(ctx)
	ue := dal.Q.User.As("ue")
	uc := dal.Q.User.As("uc")
	ci := dal.Q.Client.As("ci")
	u := dal.Q.User

	var ticketDb *model.TicketSingle
	err := tid.Select(ci.ClientID.As("client_id_str"), ti.CreatedBy, ti.Description, ue.NickName.As("editor_name"), ti.ID, ti.NumOfMultipleImage,
		ti.NumOfSingleImage, ti.Priority, uc.NickName.As("qc_name"), ti.Status, ti.Title, ti.IsActive, ti.EditorID, ti.QcID).
		Join(ci, ti.ClientID.EqCol(ci.ID)).
		Join(ue, ti.EditorID.EqCol(ue.UserID)).
		Join(uc, ti.QcID.EqCol(uc.UserID)).
		Where(ti.ID.Eq(ticketId)).Scan(&ticketDb)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrTicketNotFound
		}
		return nil, err
	}

	if !ticketDb.IsActive {
		return nil, apperror.ErrTicketHasBeenDeleted
	}

	payload := ctx.Value(model.AuthorizationPayloadKey).(*interactorinterface.Payload)
	if payload == nil {
		return nil, apperror.ErrUserNotExists
	}

	// get user type
	user, err := u.WithContext(ctx).Where(u.UserID.Eq(payload.UserID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrUserNotExists
		}
		return nil, err
	}

	if user.ContractType == string(p2mapi.FREELANCE) && payload.UserID != ticketDb.EditorID && payload.UserID != ticketDb.QcID {
		return nil, apperror.ErrViewPermissionDenied
	}

	return ticketDb.FromTicket(), err
}

func (t *TicketManagement) DeleteTicket(ctx context.Context, ticketID int64) error {
	u := dal.Q.Ticket

	_, err := u.WithContext(ctx).Where(u.ID.Eq(ticketID)).Where(u.IsActive.Is(true)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrTicketNotFound
		}
		return err
	}

	return t.txManager.TransactionExec(ctx, func(childCtx context.Context) error {
		tx := childCtx.Value(txTransactionKey).(*dal.QueryTx)

		// Delete links
		_, err = tx.Link.WithContext(childCtx).Where(tx.Link.TicketID.Eq(ticketID)).Delete()
		if err != nil {
			return err
		}

		// Delete histories
		_, err = tx.History.WithContext(childCtx).Where(tx.History.TicketID.Eq(ticketID)).Delete()
		if err != nil {
			return err
		}

		// Delete comments
		_, err = tx.Comment.WithContext(childCtx).Where(tx.Comment.TicketID.Eq(ticketID)).Delete()
		if err != nil {
			return err
		}

		_, err = tx.Ticket.WithContext(childCtx).Where(tx.Ticket.ID.Eq(ticketID)).UpdateSimple(tx.Ticket.IsActive.Value(false))
		return err
	})
}

func (t *TicketManagement) AddTicketAutoHelper(ctx context.Context, body p2mapi.CreateTicketAutoBody) error {
	// Get unassigned user to create ticket later
	u := dal.Q.User
	unassignedUser, err := u.WithContext(ctx).Where(u.NickName.Eq("unassigned")).First()
	if err != nil {
		return err
	}

	// Parse folders in body
	newTickets := make([]*model.Ticket, 0)
	for _, folder := range body.Folders {
		// parse the folder's path to get title and client_id
		// i.e., /volume5/FOR DEVELOPER/CLIENTS/SAW/UPLOAD/2024/7/21/LIBERTY BELL/LIBERTY BELL.zip

		parts := strings.Split(folder, "/")
		// TODO: We may need to update this logic when apply in PROD
		if len(parts) < 10 {
			continue
		}

		// Create client if needed
		client, err := t.clientManagement.CreateClient(ctx, p2mapi.ClientBody{
			ClientId: parts[4],
		})
		if err != nil && !errors.Is(err, apperror.ErrClientHasIDExists) {
			return err
		}

		newTickets = append(newTickets, &model.Ticket{
			ClientID:  client.Id,
			Title:     parts[9],
			CreatedBy: string(p2mapi.AUTO),
			IsActive:  true,
			Status:    string(p2mapi.BACKLOG),
			QcID:      unassignedUser.UserID,
			EditorID:  unassignedUser.UserID,
			Priority:  string(p2mapi.NORMAL),
		})
	}

	return t.txManager.TransactionExec(ctx, func(childCtx context.Context) error {
		tx := childCtx.Value(txTransactionKey).(*dal.QueryTx)

		// Save request to nas_requests table
		err := tx.NasRequest.WithContext(childCtx).Create(&model.NasRequest{
			NasID:   body.NasId,
			Payload: strings.Join(body.Folders, "\n"),
			Status:  "DONE",
		})

		if err != nil {
			return err
		}

		// Create ticket
		if len(newTickets) > 0 {
			err = tx.Ticket.WithContext(childCtx).CreateInBatches(newTickets, 50)
			if err != nil {
				return err
			}

			// Create history
			for _, ticket := range newTickets {
				err := tx.History.WithContext(childCtx).Create(&model.History{
					TicketID:    ticket.ID,
					Action:      fmt.Sprintf("Ticket is created by %s", string(p2mapi.AUTO)),
					PerformedBy: unassignedUser.UserID,
				})
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (t *TicketManagement) AddTicketAuto(ctx context.Context, body p2mapi.CreateTicketAutoBody) error {
	err := t.AddTicketAutoHelper(ctx, body)
	if err != nil {
		dal.Q.NasRequest.WithContext(ctx).Create(&model.NasRequest{
			NasID:   body.NasId,
			Payload: strings.Join(body.Folders, "\n"),
			Status:  "FAILED",
			Error:   err.Error(),
		})
	}
	return err
}
