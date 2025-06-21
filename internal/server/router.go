package server

import (
	"multi-tenant-messaging-app/internal/config"
	"multi-tenant-messaging-app/internal/handler"
	"multi-tenant-messaging-app/internal/repository"
	"multi-tenant-messaging-app/internal/service"

	_ "multi-tenant-messaging-app/docs" // Import generated docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Multi-Tenant Messaging App API
// @version 1.0
// @description This is the API documentation for the Multi-Tenant Messaging App.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Swagger (optional, if Swagger documentation is available)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize repository
	messageRepository := repository.NewMessageRepository(cfg.GormDB)

	// Initialize services
	messageService := service.NewMessageService(messageRepository, cfg.RabbitMQConn)
	tenantService := service.NewTenantService(messageRepository, cfg.RabbitMQConn)

	// Initialize handlers
	messageHandler := handler.NewMessageHandler(messageService)
	tenantHandler := handler.NewTenantHandler(tenantService)

	// Define routes
	api := r.Group("/api")
	{
		api.POST("/tenants", tenantHandler.CreateTenantHandler)
		api.DELETE("/tenants/:id", tenantHandler.DeleteTenantHandler)
		api.PUT("/tenants/:id/config/concurrency", tenantHandler.UpdateConcurrencyHandler)

		api.GET("/tenants/:tenant_id/messages", messageHandler.FetchMessages)
	}

	return r
}
