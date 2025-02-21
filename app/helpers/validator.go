package helpers

import (
	"fmt"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"

	"github.com/aadarshvelu/bms/app/repo"
)

// BookPostRequestPayloadValidator validates all fields for new book creation
func BookPostRequestPayloadValidator(c *gin.Context, payload repo.CreateBookPayload) bool {
	// Validate Title
	if !utf8.ValidString(payload.Title) || len(payload.Title) == 0 || len(payload.Title) > 512 {
		OnRequestError(c, "Title must be between 1 and 512 characters", http.StatusBadRequest)
		return false
	}

	// Validate Author
	if !utf8.ValidString(payload.Author) || len(payload.Author) == 0 || len(payload.Author) > 256 {
		OnRequestError(c, "Author name must be between 1 and 256 characters", http.StatusBadRequest)
		return false
	}

	// Validate Year
	currentYear := time.Now().Year()
	if *payload.Year <= 0 {
		OnRequestError(c, "Year is required and must be greater than 0", http.StatusBadRequest)
		return false
	}
	if *payload.Year > currentYear {
		OnRequestError(c, fmt.Sprintf("Year cannot be in the future (max: %d)", currentYear), http.StatusBadRequest)
		return false
	}

	return true
}

// UpdateBookPayloadValidator validates the update book payload
func UpdateBookPayloadValidator(c *gin.Context, book repo.UpdateBookPayload) bool {
	if book.Title != nil {
		if !utf8.ValidString(*book.Title) || len(*book.Title) == 0 || len(*book.Title) > 512 {
			OnRequestError(c, "Title must be between 1 and 512 characters", http.StatusBadRequest)
			return false
		}
	}

	fmt.Println("Error vallll:", book.Author)

	if book.Author != nil {
		if !utf8.ValidString(*book.Author) || len(*book.Author) == 0 || len(*book.Author) > 256 {
			OnRequestError(c, "Author name must be between 1 and 256 characters", http.StatusBadRequest)
			return false
		}
	}

	if book.Year != nil {
		currentYear := time.Now().Year()
		if *book.Year <= 0 || *book.Year > currentYear {
			OnRequestError(c, fmt.Sprintf("Year must be between 1 and %d", currentYear), http.StatusBadRequest)
			return false
		}
	}

	// Ensure at least one field is provided
	if book.Title == nil && book.Author == nil && book.Year == nil {
		OnRequestError(c, "At least one field must be provided for update", http.StatusBadRequest)
		return false
	}

	return true
}
