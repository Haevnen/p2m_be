package model

import (
	p2mapi "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/pkg/constants"
)

func (c *History) FromHistory() *p2mapi.HistoryResponse {
	return &p2mapi.HistoryResponse{
		Action:    c.Action,
		CreatedAt: c.CreatedAt.Format(constants.DateTimeFormat),
		Id:        c.ID,
		NickName:  c.NickName,
		TicketId:  c.TicketID,
	}
}
