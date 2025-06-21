package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	GormDB         *gorm.DB
	ServerPort     string
	JwtSecret      string
	RabbitMQURL    string
	RabbitMQConn   *amqp.Connection
	DBURL          string
	DefaultWorkers int
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	workers, _ := strconv.Atoi(os.Getenv("DEFAULT_WORKERS"))

	// Initialize GORM database connection
	dbURL := os.Getenv("DATABASE_URL")
	gormDB, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	return &Config{
		GormDB:         gormDB,
		ServerPort:     os.Getenv("SERVER_PORT"),
		JwtSecret:      os.Getenv("JWT_SECRET"),
		RabbitMQURL:    os.Getenv("RABBITMQ_URL"),
		DBURL:          dbURL,
		DefaultWorkers: workers,
	}
}
