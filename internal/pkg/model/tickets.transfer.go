package model

import (
	p2m_api "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/pkg/constants"
)

type TicketSingle struct {
	Ticket
	EditorName  string `gorm:"column:editor_name"`
	QcName      string `gorm:"column:qc_name"`
	ClientIdStr string `gorm:"column:client_id_str"`
}

func (t *TicketSingle) FromTicket() *p2m_api.SingleTicketResponse {

	return &p2m_api.SingleTicketResponse{
		ClientId:           t.ClientIdStr,
		CreatedBy:          p2m_api.CreatedBy(t.CreatedBy),
		Description:        t.Description,
		EditorName:         t.EditorName,
		Id:                 t.ID,
		NumOfMultipleImage: t.NumOfMultipleImage,
		NumOfSingleImage:   t.NumOfSingleImage,
		Priority:           p2m_api.Priority(t.Priority),
		QcName:             t.QcName,
		Status:             p2m_api.Status(t.Status),
		Title:              t.Title,
	}
}

type TicketDashboard struct {
	TicketSingle
	EditingStyle string `gorm:"column:editing_style"`
}

func (t *TicketDashboard) FromTicket() *p2m_api.DashboardResponse {
	return &p2m_api.DashboardResponse{
		ClientId:           t.ClientIdStr,
		EditingStyle:       t.EditingStyle,
		NumOfMultipleImage: t.NumOfMultipleImage,
		NumOfSingleImage:   t.NumOfSingleImage,
		Title:              t.Title,
		Date:               t.CreatedAt.Format(constants.DateTimeFormat),
		TicketId:           t.ID,
		QcName:             t.QcName,
		EditorName:         t.EditorName,
	}
}
