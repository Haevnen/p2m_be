package interactor

import (
	"context"
	"sort"
	"time"

	p2m_api "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/apperror"
	"github.com/Haevnen/p2m_be/internal/pkg/dal"
	"github.com/Haevnen/p2m_be/internal/pkg/model"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
	"github.com/Haevnen/p2m_be/pkg/util"
)

type DashboardManagement struct {
}

func NewDashboardManagement() *DashboardManagement {
	return &DashboardManagement{}
}

func (d *DashboardManagement) GetDailyDashboard(ctx context.Context, from, to time.Time, usePaging bool, page, pageSize int) ([]p2m_api.DashboardResponse, int64, error) {
	// If check from, to are UTC
	// we will convert automatically to ICT
	if util.IsUTC(from) {
		if fromServerTime, err := util.ConvertToServerTimeZone(from); err != nil {
			return nil, 0, err
		} else {
			from = *fromServerTime
		}
	}
	from = util.ChangeToUTC(from)

	if util.IsUTC(to) {
		if toServerTime, err := util.ConvertToServerTimeZone(to); err != nil {
			return nil, 0, err
		} else {
			to = *toServerTime
		}
	}
	to = util.ChangeToUTC(to)

	// Get payload to check if user is admin
	payload := ctx.Value(model.AuthorizationPayloadKey).(*interactorinterface.Payload)
	if payload == nil {
		return nil, 0, apperror.ErrUserNotExists
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
		Order(ti.CreatedAt.Asc(), ci.ClientID)

	var err error

	// Separate query to count the total records that match the filter
	var countQuery int64
	if payload.IsAdmin {
		countQuery, err = ticketQuery.Count()
		if err != nil {
			return nil, 0, err
		}
	} else {
		countQuery, err = ticketQuery.Where(tic.Where(ti.EditorID.Eq(payload.UserID)).Or(ti.QcID.Eq(payload.UserID))).Count()
		if err != nil {
			return nil, 0, err
		}
	}

	// Apply pagination if usePaging is true
	if usePaging {
		offset := (page - 1) * pageSize
		ticketQuery = ticketQuery.Offset(offset).Limit(pageSize)
	}

	var ticketDashboard []*model.TicketDashboard

	if payload.IsAdmin {
		// Get all tickets
		err = ticketQuery.Scan(&ticketDashboard)
		if err != nil {
			return nil, 0, err
		}
	} else {
		// Only get ticket relate to current user
		err = ticketQuery.Where(tic.Where(ti.EditorID.Eq(payload.UserID)).Or(ti.QcID.Eq(payload.UserID))).Scan(&ticketDashboard)
		if err != nil {
			return nil, 0, err
		}
	}

	// Try to order by date instead of date-time
	sort.Slice(ticketDashboard, func(i, j int) bool {
		// Compare year, month, and day only
		if ticketDashboard[i].CreatedAt.Year() != ticketDashboard[j].CreatedAt.Year() {
			return ticketDashboard[i].CreatedAt.Year() < ticketDashboard[j].CreatedAt.Year()
		}
		if ticketDashboard[i].CreatedAt.Month() != ticketDashboard[j].CreatedAt.Month() {
			return ticketDashboard[i].CreatedAt.Month() < ticketDashboard[j].CreatedAt.Month()
		}
		if ticketDashboard[i].CreatedAt.Day() != ticketDashboard[j].CreatedAt.Day() {
			return ticketDashboard[i].CreatedAt.Day() < ticketDashboard[j].CreatedAt.Day()
		}
		return ticketDashboard[i].ClientIdStr < ticketDashboard[j].ClientIdStr
	})

	var res []p2m_api.DashboardResponse
	for _, td := range ticketDashboard {
		res = append(res, td.FromTicket())
	}
	return res, countQuery, nil
}
