package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	
	"github.com/aadarshvelu/bms/app/services"
)

// @title Book Management System API
// @version 1.0
// @description Book Management System API with description
// @BasePath /api/v1
func SetupRoutes(r *gin.Engine) {
	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Health check
		v1.GET("/ping", services.PingHandler)

		// Book routes
		books := v1.Group("/books")
		{
			books.POST("", services.CreateBook)
			books.GET("", services.GetBooks)
			books.GET("/:id", services.GetBook)
			books.PUT("/:id", services.UpdateBook)
			books.DELETE("/:id", services.DeleteBook)
		}
	}

	// swagger
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
