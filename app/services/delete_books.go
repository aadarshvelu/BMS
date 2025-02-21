package services

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/aadarshvelu/bms/app/cache"
	"github.com/aadarshvelu/bms/app/events"
	"github.com/aadarshvelu/bms/app/helpers"
	"github.com/aadarshvelu/bms/app/repo"
	"github.com/aadarshvelu/bms/pkg/logs"
)

var logger, _ = logs.NewLogger("app.log")

// DeleteBook godoc
// @Summary Delete a book
// @Description Delete a book by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} helpers.ResponseOk
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 404 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router       /api/v1/books/{id} [delete]
func DeleteBook(c *gin.Context) {
	id, ok := helpers.ParseRequestBookId(c)
	if !ok {
		return
	}

	rowsAffected, err := repo.DeleteBook(id)

	if err != nil {
		helpers.OnRequestError(c, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		helpers.OnRequestError(c, "Book not found", http.StatusNotFound)
		return
	}

	// Invalidate both caches
	if err := cache.InvalidateBooks(); err != nil {
		logger.LogError("Failed to invalidate books cache", err)
	}

	if err := cache.InvalidateBook(id); err != nil {
		logger.LogError("Failed to invalidate book cache", err)
	}
	
	// publish an event to kafka
	if err := events.PublishBookEvent("DELETE", uint(id), nil); err != nil {
		logger.LogError("Failed to publish book deletion event", err)
	}

	helpers.OnRequestSuccess(c, nil, "Book deleted successfully", http.StatusOK)
}
