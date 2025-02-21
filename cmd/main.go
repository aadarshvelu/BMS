package main

import (
	_ "github.com/aadarshvelu/bms/cmd/docs"
	"github.com/gin-gonic/gin"

	"github.com/aadarshvelu/bms/app/handler"
	"github.com/aadarshvelu/bms/config"
    "github.com/aadarshvelu/bms/app/events"
)

// @title  Book Management System API
func main() {
	// Load configuration
	config.LoadEnv()

	// Initialize database
	config.InitDB()

	// Initialize Redis
	config.InitRedis()

    // Initialize Kafka
    events.InitKafka()

	// Initialize Gin router
	r := gin.Default()

	// Setup routes
	handler.SetupRoutes(r)

	// Run the server
	r.Run(":8080") // Listen and serve on 0.0.0.0:8080
}
