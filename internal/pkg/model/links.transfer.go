package model

import (
	p2mapi "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
)

func (l *Link) ToLink(link p2mapi.LinkBody) {
	l.TicketID = link.TicketId
	l.Link = link.Url
}

func (l *Link) FromLink() *p2mapi.LinkResponse {
	return &p2mapi.LinkResponse{
		Url:      l.Link,
		TicketId: l.TicketID,
		Id:       l.ID,
	}
}
