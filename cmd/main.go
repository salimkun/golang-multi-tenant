package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"multi-tenant-messaging-app/internal/config"
	"multi-tenant-messaging-app/internal/server"

	"github.com/streadway/amqp"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize PostgreSQL connection using db.go
	db, err := config.ConnectPostgres(cfg.DBURL)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}

	// Assign the database connection to cfg.DB
	cfg.GormDB = db

	// Initialize RabbitMQ connection
	rabbitConn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	cfg.RabbitMQConn = rabbitConn

	defer rabbitConn.Close()

	// Setup router
	router := server.SetupRouter(cfg)

	// Start server
	srvErr := make(chan error)
	go func() {
		log.Println("Server running on port:", cfg.ServerPort)
		srvErr <- router.Run(":" + cfg.ServerPort)
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("Shutting down gracefully...")
	case err := <-srvErr:
		log.Fatalf("Server error: %v", err)
	}
}
