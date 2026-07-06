package http

import (
	"net/http"
	"strconv"

	"caiyun/internal/services"
	apiresponse "caiyun/pkg/response"

	"github.com/gin-gonic/gin"
)

// ExportHandler handles export requests
type ExportHandler struct {
	exportService *services.ExportService
}

// NewExportHandler creates a new export handler
func NewExportHandler(exportService *services.ExportService) *ExportHandler {
	return &ExportHandler{exportService: exportService}
}

// ExportData handles data export requests
func (h *ExportHandler) ExportData(c *gin.Context) {
	exportType := c.Query("type")
	if exportType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type parameter is required"})
		return
	}
	format := c.Query("format")
	if format == "" {
		format = "csv"
	}
	if format != "csv" && format != "json" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format must be csv or json"})
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID, _ := userIDVal.(uint)

	filters := make(map[string]interface{})
	if accountID := c.Query("account_id"); accountID != "" {
		if id, err := strconv.Atoi(accountID); err == nil {
			filters["account_id"] = id
		}
	}
	if startDate := c.Query("start_date"); startDate != "" {
		filters["start_date"] = startDate
	}
	if endDate := c.Query("end_date"); endDate != "" {
		filters["end_date"] = endDate
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	result, err := h.exportService.ExportData(&services.ExportRequest{
		ExportType: exportType,
		Format:     format,
		Filters:    filters,
		UserID:     userID,
	})
	if err != nil {
		apiresponse.BadRequest(c, err.Error())
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, result)
}

// GetExportHistory returns export history for the user
func (h *ExportHandler) GetExportHistory(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID, _ := userIDVal.(uint)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	histories, total, err := h.exportService.GetExportHistory(userID, page, pageSize)
	if err != nil {
		apiresponse.InternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      histories,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
