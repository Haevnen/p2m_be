package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Haevnen/p2m_be/internal/pkg/registry"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
)

type historyHandler struct {
	historyManagementInteractor interactorinterface.HistoryManagementInterface
}

func newHistoryHandler(registry *registry.Registry) historyHandler {
	return historyHandler{
		historyManagementInteractor: registry.HistoryManagementInterface(),
	}
}

// Get all histories of ticket
// (GET /histories/{ticket_id})
func (h historyHandler) InternalGetHistories(c *gin.Context, ticketId int64) {
	histories, err := h.historyManagementInteractor.GetAllHistoriesByTicket(c, ticketId)
	if err != nil {
		SendError(c, "get all histories error", err)
		return
	}

	c.JSON(http.StatusOK, histories)
}
