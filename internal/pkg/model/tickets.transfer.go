package model

import p2mapi "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"

func (t *Ticket) FromTicket() *p2mapi.SingleTicketResponse {

	return &p2mapi.SingleTicketResponse{
		ClientId:           t.ClientIdStr,
		CreatedBy:          p2mapi.CreatedBy(t.CreatedBy),
		Description:        t.Description,
		EditorName:         t.EditorName,
		Id:                 t.ID,
		NumOfMultipleImage: t.NumOfMultipleImage,
		NumOfSingleImage:   t.NumOfSingleImage,
		Priority:           p2mapi.Priority(t.Priority),
		QcName:             t.QcName,
		Status:             p2mapi.Status(t.Status),
		Title:              t.Title,
	}
}
