package services

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/aadarshvelu/bms/app/cache"
	"github.com/aadarshvelu/bms/app/events"
	"github.com/aadarshvelu/bms/app/helpers"
	"github.com/aadarshvelu/bms/app/repo"
)

// CreateBook handles POST /books
// @Summary Create a new book
// @Description Create a new book in the system with the following validations:
// @Description - Title: Required, length between 1 and 512 characters
// @Description - Author: Required, length between 1 and 256 characters
// @Description - Year: Required, must not be in the future
// @Tags books
// @Accept json
// @Produce json
// @Param book body repo.CreateBookPayload true "Book details"
// @Success 201 {object} models.Book "Created book with success message"
// @Failure 400 {object} helpers.ErrorResponse "Invalid request or missing required fields"
// @Failure 500 {object} helpers.ErrorResponse "Server error"
// @Router       /api/v1/books [post]
func CreateBook(c *gin.Context) {
	var bookReq repo.CreateBookPayload

	// Get request body early and store it in context
	_ = helpers.GetRequestBody(c)

	// Parse payload
	if !helpers.ParseJsonPayload(c, &bookReq) {
		return
	}

	// Validator
	if valid := helpers.BookPostRequestPayloadValidator(c, bookReq); !valid {
		return
	}

	// Create book model from request
	record, err := repo.CreateBook(bookReq)
	if err != nil {
		helpers.OnRequestError(
			c,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	// Invalidate both caches
	if err := cache.InvalidateBooks(); err != nil {
		logger.LogError("Failed to invalidate books cache", err)
	}

	if err := events.PublishBookEvent("CREATE", record.ID, record); err != nil {
		logger.LogError("Failed to publish book creation event", err)
	}

	helpers.OnRequestSuccess(
		c,
		record,
		"Book created successfully",
		http.StatusCreated,
	)
}
