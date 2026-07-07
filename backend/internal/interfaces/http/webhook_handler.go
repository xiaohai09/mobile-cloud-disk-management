package http

import (
	"net/http"
	"strconv"

	"caiyun/internal/services"
	apiresponse "caiyun/pkg/response"

	"github.com/gin-gonic/gin"
)

// WebhookHandler handles webhook management requests
type WebhookHandler struct {
	webhookService *services.WebhookService
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(webhookService *services.WebhookService) *WebhookHandler {
	return &WebhookHandler{webhookService: webhookService}
}

// ListEndpoints lists user's webhook endpoints
func (h *WebhookHandler) ListEndpoints(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		apiresponse.Unauthorized(c, "unauthorized")
		return
	}
	uid, _ := userID.(uint)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	endpoints, total, err := h.webhookService.ListEndpoints(uid, page, pageSize)
	if err != nil {
		apiresponse.InternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"endpoints": endpoints,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CreateEndpoint creates a new webhook endpoint
func (h *WebhookHandler) CreateEndpoint(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		apiresponse.Unauthorized(c, "unauthorized")
		return
	}
	uid, _ := userID.(uint)

	var req services.CreateEndpointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(c, err.Error())
		return
	}
	req.UserID = uid

	endpoint, err := h.webhookService.CreateEndpoint(&req)
	if err != nil {
		apiresponse.BadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, endpoint)
}

// UpdateEndpoint updates a webhook endpoint
func (h *WebhookHandler) UpdateEndpoint(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		apiresponse.Unauthorized(c, "unauthorized")
		return
	}
	uid, _ := userID.(uint)

	idUint, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		apiresponse.BadRequest(c, "invalid endpoint id")
		return
	}

	var req services.UpdateEndpointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(c, err.Error())
		return
	}

	endpoint, err := h.webhookService.UpdateEndpoint(uint(idUint), uid, &req)
	if err != nil {
		apiresponse.BadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, endpoint)
}

// DeleteEndpoint deletes a webhook endpoint
func (h *WebhookHandler) DeleteEndpoint(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		apiresponse.Unauthorized(c, "unauthorized")
		return
	}
	uid, _ := userID.(uint)

	idUint, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		apiresponse.BadRequest(c, "invalid endpoint id")
		return
	}

	if err := h.webhookService.DeleteEndpoint(uint(idUint), uid); err != nil {
		apiresponse.BadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

// TestWebhook tests a webhook endpoint
func (h *WebhookHandler) TestWebhook(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		apiresponse.Unauthorized(c, "unauthorized")
		return
	}
	uid, _ := userID.(uint)

	idUint, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		apiresponse.BadRequest(c, "invalid endpoint id")
		return
	}

	if err := h.webhookService.TestEndpoint(uint(idUint), uid); err != nil {
		apiresponse.BadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "test event sent"})
}
