package handler

import (
	"github.com/Haevnen/p2m_be/internal/pkg/registry"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
	"github.com/gin-gonic/gin"
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
func (h historyHandler) InternalGetHistories(c *gin.Context, ticketId int64) {}
