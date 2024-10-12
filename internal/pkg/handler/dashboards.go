package handler

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	p2mapi "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
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

// InternalGetDailyDashboard (GET /dashboards/daily)
func (h dashboardHandler) InternalGetDailyDashboard(c *gin.Context, params p2mapi.InternalGetDailyDashboardParams) {
	dailyDashboard, err := h.dashboardInteractor.GetDailyDashboard(c, params.Date.Time, params.Date.Time)
	if err != nil {
		SendError(c, "get daily dashboard error", err)
		return
	}

	c.JSON(http.StatusOK, dailyDashboard)
}

// InternalExportDashboard (GET /dashboards/export)
func (h dashboardHandler) InternalExportDashboard(c *gin.Context, params p2mapi.InternalExportDashboardParams) {
	difference := params.StartTime.Sub(params.EndTime)
	if int64(difference.Hours()/24) > model.MaxExportRangeInDay {
		SendError(c, "export select time over range", apperror.ErrExportTimeOverRange)
		return
	}

	data, err := h.dashboardInteractor.GetDailyDashboard(c, params.StartTime, params.EndTime)
	if err != nil {
		SendError(c, "export error", err)
		return
	}

	dataConvert := make([]*model.TicketExport, len(data))
	for i, v := range data {
		d := (model.DashboardResponse)(*v)
		dataConvert[i] = d.FromDashboardResponse()
	}

	// Set the correct headers
	fileName := "dashboard_" + time.Now().Format("20060102_1504") + ".csv"
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "text/csv")

	// Create a CSV writer
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write the CSV header from struct tags
	if len(dataConvert) > 0 {
		header := getCSVHeader(*dataConvert[0])
		if err := writer.Write(header); err != nil {
			SendError(c, "writing CSV header error", err)
			return
		}
	}

	// Write the CSV rows
	for _, d := range dataConvert {
		record := getCSVRecord(*d)
		if err := writer.Write(record); err != nil {
			SendError(c, "writing CSV record error", err)
			return
		}
	}
}

// getCSVHeader returns the CSV header based on the `csv` tags in the struct
func getCSVHeader(data interface{}) []string {
	val := reflect.ValueOf(data)
	typ := val.Type()

	var header []string
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("csv")
		if tag != "" {
			header = append(header, tag)
		} else {
			header = append(header, typ.Field(i).Name) // Fallback to field name if no tag
		}
	}
	return header
}

// getCSVRecord returns a slice of string values for each struct field
func getCSVRecord(data interface{}) []string {
	val := reflect.ValueOf(data)

	var record []string
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		switch field.Kind() {
		case reflect.String:
			record = append(record, field.String())
		case reflect.Int:
			record = append(record, strconv.Itoa(int(field.Int())))
		default:
			record = append(record, fmt.Sprintf("%v", field.Interface()))
		}
	}
	return record
}
