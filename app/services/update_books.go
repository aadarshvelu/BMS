package services

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/aadarshvelu/bms/app/cache"
	"github.com/aadarshvelu/bms/app/events"
	"github.com/aadarshvelu/bms/app/helpers"
	"github.com/aadarshvelu/bms/app/repo"
)

// UpdateBook handles PUT /books/:id
// @Summary Update a book
// @Description Update an existing book in the system
// @Description - Title length should be between 1 and 512 characters
// @Description - Author length should be between 1 and 256 characters
// @Description - Year must not be in the future
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body repo.UpdateBookPayload true "Book details"
// @Success 200 {object} helpers.ResponseOk
// @Failure 400 {object} helpers.ErrorResponse "Invalid request or validation error"
// @Failure 404 {object} helpers.ErrorResponse "Book not found"
// @Failure 500 {object} helpers.ErrorResponse "Server error"
// @Router       /api/v1/books/{id} [put]
func UpdateBook(c *gin.Context) {
	// Get request body early and store it in context
	_ = helpers.GetRequestBody(c)

	// Parse ID
	id, ok := helpers.ParseRequestBookId(c)
	if !ok {
		return
	}

	// Parse request body
	var bookReq repo.UpdateBookPayload
	if !helpers.ParseJsonPayload(c, &bookReq) {
		return
	}

	// Validate request
	if valid := helpers.UpdateBookPayloadValidator(c, bookReq); !valid {
		return
	}

	// Update book
	record, err := repo.UpdateBook(id, bookReq)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "record not found" {
			statusCode = http.StatusNotFound
		}

		helpers.OnRequestError(c, err.Error(), statusCode)
		return
	}

	// Invalidate both caches
	if err := cache.InvalidateBooks(); err != nil {
		logger.LogError("Failed to invalidate books cache", err)
	}

	if err := cache.InvalidateBook(id); err != nil {
		logger.LogError("Failed to invalidate book cache", err)
	}

	if err := events.PublishBookEvent("UPDATE", record.ID, record); err != nil {
		logger.LogError("Failed to publish book update event", err)
	}

	helpers.OnRequestSuccess(c, record, "Book updated successfully", http.StatusOK)
}
