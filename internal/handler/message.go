package handler

import (
	"multi-tenant-messaging-app/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

func (h *MessageHandler) FetchMessages(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	cursor := c.Query("cursor")
	limit := 10 // Default limit, can be adjusted or made configurable

	messages, lastID, err := h.messageService.FetchMessages(tenantID, cursor, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"last_id":  lastID,
	})
}
