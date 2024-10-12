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
	EditingStyle       string `gorm:"column:editing_style"`
	EditorContractType string `gorm:"column:editor_contract_type"`
	QcContractType     string `gorm:"column:qc_contract_type"`
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
		QcContractType:     p2m_api.ContractType(t.QcContractType),
		EditorName:         t.EditorName,
		EditorContractType: p2m_api.ContractType(t.EditorContractType),
	}
}

type TicketExport struct {
	ClientId           string `csv:"client_id"`
	TicketId           int64  `csv:"task"`
	Title              string `csv:"title"`
	EditingStyle       string `csv:"editing_style"`
	NumOfSingleImage   int32  `csv:"num_of_single_image"`
	NumOfMultipleImage int32  `csv:"num_of_multiple_image"`
	QcName             string `csv:"qc_name"`
	QcContractType     string `csv:"qc_contract_type"`
	EditorName         string `csv:"editor_name"`
	EditorContractType string `csv:"editor_contract_type"`
}

type DashboardResponse p2m_api.DashboardResponse

func (t DashboardResponse) FromDashboardResponse() *TicketExport {
	return &TicketExport{
		ClientId:           t.ClientId,
		TicketId:           t.TicketId,
		Title:              t.Title,
		EditingStyle:       t.EditingStyle,
		NumOfMultipleImage: t.NumOfMultipleImage,
		NumOfSingleImage:   t.NumOfSingleImage,
		QcName:             t.QcName,
		QcContractType:     string(t.QcContractType),
		EditorName:         t.EditorName,
		EditorContractType: string(t.EditorContractType),
	}
}
