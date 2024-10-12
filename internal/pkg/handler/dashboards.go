package handler

import (
	"net/http"

	p2m_api "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/pkg/registry"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
	"github.com/gin-gonic/gin"
)

type dashboardHandler struct {
	dashboardInteractor interactorinterface.DashboardInterface
}

func newDashboardHandler(registry *registry.Registry) dashboardHandler {
	return dashboardHandler{
		dashboardInteractor: registry.DashboardInterface(),
	}
}

// (GET /dashboards/daily)
func (h dashboardHandler) InternalGetDailyDashboard(c *gin.Context, params p2m_api.InternalGetDailyDashboardParams) {
	dailyDashboard, err := h.dashboardInteractor.GetDailyDashboard(c, params.Date.Time)
	if err != nil {
		SendError(c, "get daily dashboard error", err)
		return
	}

	c.JSON(http.StatusOK, dailyDashboard)
}

// (GET /dashboards/export)
func (h dashboardHandler) InternalExportDashboard(c *gin.Context, params p2m_api.InternalExportDashboardParams) {

}
