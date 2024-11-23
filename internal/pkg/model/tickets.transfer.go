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

func (t *TicketDashboard) FromTicket() p2m_api.DashboardResponse {
	return p2m_api.DashboardResponse{
		ClientId:           t.ClientIdStr,
		EditingStyle:       t.EditingStyle,
		NumOfMultipleImage: t.NumOfMultipleImage,
		NumOfSingleImage:   t.NumOfSingleImage,
		Title:              t.Title,
		Date:               t.CreatedAt.Format(constants.DateFormat),
		TicketId:           t.ID,
		QcName:             t.QcName,
		QcContractType:     p2m_api.ContractType(t.QcContractType),
		EditorName:         t.EditorName,
		EditorContractType: p2m_api.ContractType(t.EditorContractType),
	}
}

type TicketExport struct {
	Date               string `csv:"Date"`
	ClientId           string `csv:"Client"`
	TicketId           int64  `csv:"Task"`
	Title              string `csv:"Title of ticket"`
	EditingStyle       string `csv:"Editing style"`
	NumOfSingleImage   int32  `csv:"Quantity single"`
	NumOfMultipleImage int32  `csv:"Quantity multiple"`
	EditorName         string `csv:"Editor"`
	EditorContractType string `csv:"Editor contract type"`
	QcName             string `csv:"QC"`
	QcContractType     string `csv:"QC contract type"`
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
		Date:               t.Date,
	}
}
