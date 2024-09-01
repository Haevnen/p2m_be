package handler

import (
	"net/http"

	p2m_api "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/apperror"
	"github.com/Haevnen/p2m_be/internal/pkg/registry"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
	"github.com/gin-gonic/gin"
)

type ticketHandler struct {
	ticketManagementInteractor interactorinterface.TicketManagementInterface
}

func newTicketHandler(registry *registry.Registry) ticketHandler {
	return ticketHandler{
		ticketManagementInteractor: registry.TicketManagementInterface(),
	}
}

// Get all tickets
// (GET /tickets)
func (h ticketHandler) InternalGetTickets(c *gin.Context) {
	// TODO:
	// * Need to get payload to check if user is admin (admin user will get all tickets, other member
	// 	 only get ticket related to that member)
	// * Ticket will group by status and sort by priority desc and updated_at
	// desc
	// * We query today for DONE status, all day for other status
}

// Add new ticket
// (POST /tickets/add)
func (h ticketHandler) InternalAddTicket(c *gin.Context, params p2m_api.InternalAddTicketParams) {
	var creatTicketBody p2m_api.CreateTicketBody

	if err := bindRequestBody(c, &creatTicketBody); err != nil {
		SendError(c, err.Error(), apperror.ErrInvalidRequestInput)
		return
	}

	err := h.ticketManagementInteractor.AddTicketManual(c, creatTicketBody)
	if err != nil {
		SendError(c, "add ticket error", err)
		return
	}

	c.JSON(http.StatusOK, "add ticket successfully")
}

// Remove ticket by id
// (DELETE /tickets/remove/{ticket_id})
func (h ticketHandler) InternalRemoveTicket(c *gin.Context, ticketId int, params p2m_api.InternalRemoveTicketParams) {
	// TODO:
	// * Soft delete
	// * Remove all related links, histories and comments
}

// Update ticket by id
// (PUT /tickets/update/{ticket_id})
func (h ticketHandler) InternalUpdateTicket(c *gin.Context, ticketId int64, params p2m_api.InternalUpdateTicketParams) {
	var updateTicketBody p2m_api.UpdateTicketBody
	if err := bindRequestBody(c, &updateTicketBody); err != nil {
		SendError(c, err.Error(), apperror.ErrInvalidRequestInput)
		return
	}

	err := h.ticketManagementInteractor.UpdateTicket(c, ticketId, updateTicketBody)
	if err != nil {
		SendError(c, "update ticket error", err)
		return
	}

	c.JSON(http.StatusOK, "update ticket successfully")
}

// Get ticket by id
// (GET /tickets/{ticket_id})
func (h ticketHandler) InternalGetTicket(c *gin.Context, ticketId int) {
	// TODO:
	// * We not only need to get data from ticket but also the other data like
	// (comments, links, histories and client_info) -> Does BE handle or FE will
	// send request?
}
