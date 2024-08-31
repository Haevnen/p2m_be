package model

import (
	p2mapi "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/pkg/constants"
)

func (c *History) FromHistory() *p2mapi.HistoryResponse {
	createdAt := c.CreatedAt.Format(constants.DateTimeFormat)
	return &p2mapi.HistoryResponse{
		Action:    c.Action,
		CreatedAt: &createdAt,
		Id:        c.ID,
		NickName:  c.PerformedBy,
		TicketId:  c.TicketID,
	}
}
