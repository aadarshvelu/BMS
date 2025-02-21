package helpers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/aadarshvelu/bms/pkg/logs"
)

type ResponseOk struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
	Status       int    `json:"status"`
}

var logger *logs.Logger

func init() {
	var err error
	logger, err = logs.NewLogger("api_logs.txt")
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
}

func OnRequestSuccess(c *gin.Context, data interface{}, msg string, status int) {
	response := ResponseOk{
		Data:    data,
		Message: msg,
		Status:  status,
	}

	// Log the request
	logger.LogRequest(
		c.Request.URL.String(),
		c.Request.Method,
		status,
		GetRequestBody(c),
		response,
	)

	c.JSON(status, response)
}

func OnRequestError(c *gin.Context, msg string, status int) {
	response := ErrorResponse{
		ErrorMessage: msg,
		Status:       status,
	}

	// Log the request
	logger.LogRequest(
		c.Request.URL.String(),
		c.Request.Method,
		status,
		GetRequestBody(c),
		response,
	)

	c.JSON(status, response)
}

func GetRequestBody(c *gin.Context) string {
	// Check if we already stored the body in the context
	if body, exists := c.Get("cachedRequestBody"); exists {
		return body.(string)
	}

	if c.Request.Method == "POST" || c.Request.Method == "PUT" {
		if c.Request.ContentLength == 0 {
			return "nil"
		}

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println("Error reading request body:", err)
			return "nil"
		}

		// Restore the body for subsequent reads
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Store the body in the context for reuse
		bodyStr := string(bodyBytes)
		c.Set("cachedRequestBody", bodyStr)

		return bodyStr
	}

	return "nil"
}

// ParseJsonPayload parses JSON request body into the provided struct
// Returns true if parsing is successful, false otherwise
func ParseJsonPayload(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		errMsg := err.Error()

		// Handle unmarshal type errors
		if strings.Contains(errMsg, "unmarshal") {
			// Extract field name
			parts := strings.Split(errMsg, "field ")
			if len(parts) > 1 {
				fieldName := strings.Split(parts[1], " ")[0]
				// Convert to lowercase and remove struct name if present
				fieldName = strings.ToLower(strings.Split(fieldName, ".")[len(strings.Split(fieldName, "."))-1])
				errMsg = fmt.Sprintf("Invalid %s", fieldName)
			}
		} else if strings.Contains(errMsg, "json:") {
			// Handle other JSON validation errors
			parts := strings.Split(errMsg, "json:")
			if len(parts) > 1 {
				errMsg = strings.TrimSpace(parts[1])
			}
		}

		OnRequestError(
			c,
			errMsg,
			http.StatusBadRequest,
		)
		return false
	}
	return true
}

// ParseRequestBookId parses and validates the book ID from the request parameters
// Returns the parsed ID if successful, and a boolean indicating success/failure
func ParseRequestBookId(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		OnRequestError(
			c,
			"Invalid ID",
			http.StatusBadRequest,
		)
		return 0, false
	}
	return id, true
}
