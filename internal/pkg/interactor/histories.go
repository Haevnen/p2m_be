package interactor

import (
	"context"

	p2mapi "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/pkg/dal"
	"github.com/Haevnen/p2m_be/internal/pkg/model"
)

type HistoryManagement struct {
}

func NewHistoryManagement() *HistoryManagement {
	return &HistoryManagement{}
}

func (ci *HistoryManagement) GetAllHistoriesByTicket(ctx context.Context, ticketID int64) ([]*p2mapi.HistoryResponse, error) {
	c := dal.Q.History
	u := dal.Q.User

	var histories []*model.HistoryWithName
	err := c.WithContext(ctx).Select(c.ID, c.TicketID, c.Action, c.PerformedBy, c.CreatedAt, u.NickName).
		Join(u, c.PerformedBy.EqCol(u.UserID)).
		Where(c.TicketID.Eq(ticketID)).Order(c.CreatedAt.Desc()).Scan(&histories)

	if err != nil {
		return nil, err
	}

	res := make([]*p2mapi.HistoryResponse, 0, len(histories))
	for _, history := range histories {
		res = append(res, history.FromHistory())
	}

	return res, nil
}
