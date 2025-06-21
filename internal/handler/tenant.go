package handler

import (
	"multi-tenant-messaging-app/internal/payload"
	"multi-tenant-messaging-app/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TenantHandler struct {
	tenantService *service.TenantService
}

func NewTenantHandler(tenantService *service.TenantService) *TenantHandler {
	return &TenantHandler{tenantService: tenantService}
}

// CreateTenantHandler godoc
// @Summary Create a new tenant and publish a message
// @Description Creates a new tenant and starts a consumer for the tenant. Also publishes an initial message to the tenant's queue.
// @Tags Tenants
// @Accept json
// @Produce json
// @Param tenant body payload.TenantRequest true "Tenant ID and initial payload"
// @Success 200 {object} map[string]string "message: Tenant created and message published"
// @Failure 400 {object} map[string]string "error: Invalid request"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Router /api/tenants [post]
func (h *TenantHandler) CreateTenantHandler(c *gin.Context) {
	var req payload.TenantRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if _, err := uuid.Parse(req.TenantID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	// Mulai konsumen untuk tenant
	if err := h.tenantService.StartConsumer(req.TenantID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Kirim pesan ke RabbitMQ queue
	if err := h.tenantService.PublishToTenantQueue(req.TenantID, req.Payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish message: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenant created"})
}

// DeleteTenantHandler godoc
// @Summary Delete a tenant and stop its consumer
// @Description Deletes a tenant and stops the associated RabbitMQ consumer.
// @Tags Tenants
// @Param id path string true "Tenant ID"
// @Success 200 {object} map[string]string "message: Tenant deleted"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Router /api/tenants/{id} [delete]
func (h *TenantHandler) DeleteTenantHandler(c *gin.Context) {
	tenantID := c.Param("id")

	if err := h.tenantService.StopConsumer(tenantID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenant removed"})
}

// UpdateConcurrencyHandler godoc
// @Summary Update worker concurrency for a tenant
// @Description Updates the number of workers processing messages for a specific tenant.
// @Tags Tenants
// @Param id path string true "Tenant ID"
// @Param workers body payload.UpdateConcurrencyRequest true "Number of workers"
// @Success 200 {object} map[string]string "message: Concurrency updated"
// @Failure 400 {object} map[string]string "error: Invalid workers configuration"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Router /api/tenants/{id}/config/concurrency [put]
func (h *TenantHandler) UpdateConcurrencyHandler(c *gin.Context) {
	tenantID := c.Param("id")
	if _, err := uuid.Parse(tenantID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant UUID"})
		return
	}

	var req struct {
		Workers int `json:"workers"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Workers <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workers configuration"})
		return
	}

	if err := h.tenantService.UpdateWorkerCount(tenantID, req.Workers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Concurrency updated"})
}
