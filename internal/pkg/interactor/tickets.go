package interactor

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
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
	"github.com/Haevnen/p2m_be/pkg/logger"
	"github.com/Haevnen/p2m_be/pkg/util"
)

const (
	defaultDescription = `<p class="editor-paragraph"><br></p>`
)

var exceptionClientID = map[string]bool{
	"COR": true,
	"DAB": true,
	"BRH": true,
}

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
					ClientID: newTicket.ClientID,
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
		ti := dal.Q.Ticket
		_, err = tx.Ticket.WithContext(childCtx).Select(ti.ID, ti.Title, ti.Status, ti.QcID, ti.EditorID, ti.Priority, ti.ClientID, ti.Description, ti.CreatedBy, ti.NumOfMultipleImage, ti.NumOfSingleImage, ti.IsActive, ti.CreatedAt).
			Where(tx.Ticket.ID.Eq(ticketID)).Where(tx.Ticket.IsActive.Is(true)).Omit(tx.Ticket.UpdatedAt).Updates(&ticket)
		return err
	})
}

func (t *TicketManagement) GetAllTicketsByContractType(ctx context.Context) ([]*p2mapi.ListTicketItem, error) {
	ti := dal.Q.Ticket
	tid := ti.WithContext(ctx)
	u := dal.Q.User
	ci := dal.Q.Client.As("ci")

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
	ticketQuery := tid.Join(ci, ti.ClientID.EqCol(ci.ID)).
		Select(ti.Title, ti.ID, ti.Status, ti.Priority, ti.QcID, ti.EditorID, ti.UpdatedAt, ti.ClientID, ti.CreatedAt).
		// Ignore deleted ticket
		Where(ti.IsActive.Is(true)).
		// List all ticket not completed (for all day)
		Where(tid.Where(ti.Status.Neq(string(p2mapi.DONE))).
			// Or List all ticket in status DONE (for current day)
			Or(tid.Where(ti.UpdatedAt.Between(util.Begin(now), util.End(now))).Where(ti.Status.Eq(string(p2mapi.DONE)))))
	// List by priority and updated_at asc
	ticketQuery.Order(ti.CreatedAt.Asc()).Order(ci.ClientID)

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

		// Get all clients
		clients, err := t.clientManagement.GetAllClient(ctx, swag.Bool(true))
		if err != nil {
			return nil, err
		}
		clientIdMapping := make(map[int32]*p2mapi.ClientResponse)
		for _, client := range clients {
			clientIdMapping[client.Id] = client
		}

		// Try to order by date instead of date-time
		sort.Slice(ticketsDb, func(i, j int) bool {
			// Compare year, month, and day only
			if ticketsDb[i].CreatedAt.Year() != ticketsDb[j].CreatedAt.Year() {
				return ticketsDb[i].CreatedAt.Year() < ticketsDb[j].CreatedAt.Year()
			}
			if ticketsDb[i].CreatedAt.Month() != ticketsDb[j].CreatedAt.Month() {
				return ticketsDb[i].CreatedAt.Month() < ticketsDb[j].CreatedAt.Month()
			}
			if ticketsDb[i].CreatedAt.Day() != ticketsDb[j].CreatedAt.Day() {
				return ticketsDb[i].CreatedAt.Day() < ticketsDb[j].CreatedAt.Day()
			}

			// If year, month, and day are the same, compare by client
			if ticketsDb[i].ClientID != ticketsDb[j].ClientID {
				return clientIdMapping[ticketsDb[i].ClientID].ClientId < clientIdMapping[ticketsDb[j].ClientID].ClientId
			}

			// If year, month, day, and client are the same, compare by ticket
			// name
			return ticketsDb[i].Title < ticketsDb[j].Title
		})

		for _, ticket := range ticketsDb {
			ticketItem := p2mapi.ListTicket{
				EditorName: userIdMapping[ticket.EditorID].NickName,
				Id:         ticket.ID,
				Priority:   p2mapi.Priority(ticket.Priority),
				QcName:     userIdMapping[ticket.QcID].NickName,
				Title:      ticket.Title,
				UpdatedAt:  ticket.UpdatedAt.Format(constants.DateTimeFormat),
				ClientId:   clientIdMapping[ticket.ClientID].ClientId,
				CreatedAt:  ticket.CreatedAt.Format(constants.DateFormat),
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

// Get unassigned user to create ticket later
func (t *TicketManagement) getUnassignedUser(ctx context.Context) (*model.User, error) {
	u := dal.Q.User
	unassignedUser, err := u.WithContext(ctx).Where(u.NickName.Eq("unassigned")).First()
	if err != nil {
		return nil, err
	}
	return unassignedUser, nil
}

// Parse folder path to get internal path
func (t *TicketManagement) parseFolderPathToGetTicketMetadata(root, folder string) (string, string, string, error) {
	indexOfRoot := strings.Index(folder, root)
	if indexOfRoot == -1 {
		return "", "", "", fmt.Errorf("invalid folder path: %s - does not start with root path", folder)
	}

	// This regex looks for any characters between / and /UPLOAD/
	re := regexp.MustCompile(`/([^/]+)/UPLOAD/`)
	matches := re.FindStringSubmatch(folder)

	var clientID = ""
	if len(matches) >= 2 {
		clientID = matches[1]
	}

	if clientID == "" {
		return "", "", "", fmt.Errorf("invalid folder path: %s - no client ID found", folder)
	}

	parts := strings.Split(folder, "/")
	if exceptionClientID[clientID] {
		// Exception path:
		// /volume3/FILES STATION/CLIENTS/COR/UPLOAD/01-07-2024/1513 Paddington
		// /volume3/BRH/UPLOAD/2024/November/1/11-1 Melinda Brad 137
		if len(parts) != 8 {
			return "", "", "", fmt.Errorf("invalid folder path: %s - wrong number of parts", folder)
		}
	} else {
		// Normal path:
		// /volume3/FILES STATION/CLIENTS/SAW/UPLOAD/2024/9/14/TestAuto
		if len(parts) != 10 {
			return "", "", "", fmt.Errorf("invalid folder path: %s - wrong number of parts", folder)
		}
	}

	// ClientID, Title, InternalLink
	return clientID, parts[len(parts)-1], folder[indexOfRoot+len(root):], nil
}

// Check if a ticket already exists for the given title and client ID
func ticketExists(ctx context.Context, expectedInternalLink string, client int32) (bool, error) {
	l := dal.Q.Link

	count, err := l.WithContext(ctx).Where(l.ClientID.Eq(client)).Where(l.Link.Eq(expectedInternalLink)).Count()
	return count > 0, err
}

func (t *TicketManagement) AddTicketAutoHelper(ctx context.Context, body p2mapi.CreateTicketAutoBody) error {
	na := dal.Q.NasServer

	unassignedUser, err := t.getUnassignedUser(ctx)
	if err != nil {
		return err
	}

	nasServer, err := na.WithContext(ctx).Where(na.NasID.Eq(body.NasId)).First()
	if err != nil {
		return err
	}

	newTickets := make([]*model.Ticket, 0)
	visitedTitles := make(map[string]bool)
	for _, folder := range body.Folders {
		clientID, title, internalLink, err := t.parseFolderPathToGetTicketMetadata(nasServer.RootPath, folder)
		logger.Infof("clientID: %s, title: %s, internalLink: %s", clientID, title, internalLink)
		if err != nil {
			logger.Error(err.Error())
			continue // Skip invalid folders path
		}

		if visitedTitles[title] {
			logger.Infof("Duplicate title: %s", title)
			continue // Skip if title already exists
		}
		visitedTitles[title] = true

		// Only create ticket if client exist
		if clientID == "" {
			logger.Errorf("Client not found: %s", clientID)
			continue
		}

		client, err := t.clientManagement.GetSingleClient(ctx, clientID)
		if err != nil && !errors.Is(err, apperror.ErrRecordNotFound) {
			logger.Errorf("Failed to get client: %s", clientID)
			return err
		}
		if client == nil {
			logger.Errorf("Client not found: %s", clientID)
			continue
		}

		// Check if ticket already exists
		expectedInternalLink := `\\` + nasServer.InternalPath + strings.ReplaceAll(internalLink, "/", `\`)
		exists, err := ticketExists(ctx, expectedInternalLink, client.Id)
		if err != nil {
			logger.Errorf("Failed to check if ticket exists: %s", title)
			return err
		}

		if exists {
			logger.Infof("Ticket already exists: %s", expectedInternalLink)
			continue // Skip if ticket already exists
		}

		newTickets = append(newTickets, &model.Ticket{
			ClientID:     client.Id,
			Title:        title,
			CreatedBy:    string(p2mapi.AUTO),
			IsActive:     true,
			Status:       string(p2mapi.BACKLOG),
			QcID:         unassignedUser.UserID,
			EditorID:     unassignedUser.UserID,
			Priority:     string(p2mapi.NORMAL),
			Description:  defaultDescription,
			InternalLink: internalLink,
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

			// Create history and link
			for _, ticket := range newTickets {
				// history
				err := tx.History.WithContext(childCtx).Create(&model.History{
					TicketID:    ticket.ID,
					Action:      fmt.Sprintf("Ticket is created by %s", string(p2mapi.AUTO)),
					PerformedBy: unassignedUser.UserID,
				})
				if err != nil {
					return err
				}
				// windows explorer link
				err = tx.Link.WithContext(childCtx).Create(&model.Link{
					TicketID: ticket.ID,
					Link:     `\\` + nasServer.InternalPath + strings.ReplaceAll(ticket.InternalLink, "/", `\`),
					ClientID: ticket.ClientID,
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
