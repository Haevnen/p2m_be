package interactor

import (
	"context"
	"time"

	p2m_api "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/apperror"
	"github.com/Haevnen/p2m_be/internal/pkg/dal"
	"github.com/Haevnen/p2m_be/internal/pkg/model"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
)

type DashboardManagement struct {
}

func NewDashboardManagement() *DashboardManagement {
	return &DashboardManagement{}
}

func (d *DashboardManagement) GetDailyDashboard(ctx context.Context, from, to time.Time) ([]*p2m_api.DashboardResponse, error) {
	// Get payload to check if user is admin
	payload := ctx.Value(model.AuthorizationPayloadKey).(*interactorinterface.Payload)
	if payload == nil {
		return nil, apperror.ErrUserNotExists
	}

	ti := dal.Q.Ticket
	tic := ti.WithContext(ctx)
	ci := dal.Q.Client.As("ci")
	ue := dal.Q.User.As("ue")
	uc := dal.Q.User.As("uc")

	ticketQuery := tic.Select(
		ci.ClientID.As("client_id_str"),
		ci.EditingStyle.As("editing_style"),
		ue.NickName.As("editor_name"),
		ue.ContractType.As("editor_contract_type"),
		uc.NickName.As("qc_name"),
		uc.ContractType.As("qc_contract_type"),
		ti.CreatedAt,
		ti.NumOfMultipleImage,
		ti.NumOfSingleImage,
		ti.Title,
		ti.ID).
		Join(ci, ti.ClientID.EqCol(ci.ID)).
		Join(ue, ti.EditorID.EqCol(ue.UserID)).
		Join(uc, ti.QcID.EqCol(uc.UserID)).
		Where(ti.IsActive.Is(true), ti.CreatedAt.Between(from, to)).
		Order(ci.ClientID.Asc(), ti.CreatedAt.Asc(), ti.Priority.Desc())

	var ticketDashboard []*model.TicketDashboard
	var err error
	if payload.IsAdmin {
		// Get all tickets
		err = ticketQuery.Scan(&ticketDashboard)
		if err != nil {
			return nil, err
		}
	} else {
		// Only get ticket relate to current user
		err = ticketQuery.Where(tic.Where(ti.EditorID.Eq(payload.UserID)).Or(ti.QcID.Eq(payload.UserID))).Scan(&ticketDashboard)
		if err != nil {
			return nil, err
		}
	}

	var res []*p2m_api.DashboardResponse
	for _, td := range ticketDashboard {
		res = append(res, td.FromTicket())
	}
	return res, nil
}
