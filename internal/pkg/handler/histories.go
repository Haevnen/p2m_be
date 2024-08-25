package handler

import (
	p2m_api "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
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

// Create new history
// (POST /histories)
func (h historyHandler) InternalCreateHistory(c *gin.Context, params p2m_api.InternalCreateHistoryParams) {
}

// Get all histories of ticket
// (GET /histories/{ticket_id})
func (h historyHandler) InternalGetHistories(c *gin.Context, ticketId int64) {}
