package handler

import (
	"multi-tenant-messaging-app/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageService service.MessageServiceInterface
}

func NewMessageHandler(messageService service.MessageServiceInterface) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

// FetchMessages godoc
// @Summary Fetch messages for a tenant
// @Description Retrieves a list of messages for a specific tenant, with optional cursor-based pagination.
// @Tags Messages
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param cursor query string false "Cursor for pagination"
// @Success 200 {object} map[string]interface{} "messages: List of messages, last_id: ID of the last message"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Router /api/tenants/{tenant_id}/messages [get]
func (h *MessageHandler) FetchMessages(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	cursor := c.Query("cursor")
	limit := 10 // Default limit, can be adjusted or made configurable

	messages, cursor, err := h.messageService.FetchMessages(tenantID, cursor, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   messages,
		"cursor": cursor,
	})
}
