package handler

import (
	"github.com/gin-gonic/gin"

	p2m_api "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/apperror"
	"github.com/Haevnen/p2m_be/internal/pkg/model"
	"github.com/Haevnen/p2m_be/internal/pkg/registry"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
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

}

// Get timerange dashboard
// (GET /dashboards/export)
func (h dashboardHandler) InternalExportDashboard(c *gin.Context, params p2m_api.InternalExportDashboardParams) {
	difference := params.StartTime.Sub(params.EndTime)
	if int64(difference.Hours()/24) > model.MaxExportRangeInDay {
		SendError(c, "export select time over range", apperror.ErrExportTimeOverRange)
		return
	}

}
