package services

import (
	"net/http"

	"github.com/aadarshvelu/bms/app/helpers"
	"github.com/gin-gonic/gin"
)

// @Summary Health check endpoint
// @Description Returns a simple message indicating the system is running
// @Tags system
// @Produce json
// @Success 200 {object}  helpers.ResponseOk "Success response with message"
// @Router       /api/v1/ping [get]
func PingHandler(c *gin.Context) {
	helpers.OnRequestSuccess(c, nil, "System is up and running!", http.StatusAccepted)
}
