package services

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/aadarshvelu/bms/app/cache"
	"github.com/aadarshvelu/bms/app/helpers"
	"github.com/aadarshvelu/bms/app/models"
	"github.com/aadarshvelu/bms/app/repo"
)

type Pagination struct {
	PrevPage        *int `json:"prevPage"`
	CurrentPage     int  `json:"currentPage"`
	NextPage        *int `json:"nextPage"`
	CurrentPagesize int  `json:"currentPagesize"`
	TotalPages      int  `json:"totalPages"`
	TotalRecords    int  `json:"totalRecords"`
}

type BooksPaginationResponse struct {
	Books      []models.Book `json:"books"`
	Pagination Pagination    `json:"pagination"`
}

// GetBooks godoc
// @Summary      Get all books
// @Description  Retrieves a list of all books from the system
// @Tags         books
// @Produce      json
// @Param        page 	   query    int     false  "Page number (default 1)"
// @Param        pagesize  query    int     false  "Number of books per page (default 10)"
// @Success      200  {object}  BooksPaginationResponse
// @Failure      400  {object}  helpers.ErrorResponse "Invalid parameters"
// @Failure      500  {object}  helpers.ErrorResponse "Server error"
// @Router       /api/v1/books [get]
func GetBooks(c *gin.Context) {
	// Parse and validate limit
	limit := 10
	if limitStr := c.DefaultQuery("pagesize", "10"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit < 1 {
			helpers.OnRequestError(c, "Invalid pagesize parameter", http.StatusBadRequest)
			return
		}

		limit = parsedLimit
	}

	// Parse and validate offset
	offset := 1
	if offsetStr := c.DefaultQuery("page", "1"); offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil || parsedOffset < 1 {
			helpers.OnRequestError(c, "Invalid page parameter", http.StatusBadRequest)
			return
		}
		offset = parsedOffset
	}

	// Try to get from cache first
	books, totalBooks, err := cache.GetBooks(limit, offset)

	if err != nil {
		// Get all books from DB if not in cache
		allBooks, err := repo.GetAllBooks()
		if err != nil {
			helpers.OnRequestError(c, err.Error(), http.StatusInternalServerError)
			return
		}

		totalBooks = len(allBooks)

		// Store all books in cache for next time
		if err := cache.SetBooks(allBooks); err != nil {
			log.Printf("Failed to cache books: %v", err)
		}

		// Calculate pagination from all books
		startIndex := (offset - 1) * limit
		endIndex := startIndex + limit

		// Validate indices
		if startIndex >= totalBooks {
			books = []models.Book{}
		} else {
			if endIndex > totalBooks {
				endIndex = len(allBooks)
			}
			books = allBooks[startIndex:endIndex]
		}
	}
	pagination := Pagination{
		CurrentPage:     offset,
		CurrentPagesize: limit,
		TotalRecords:    totalBooks,
		TotalPages:      (totalBooks + limit - 1) / limit,
	}

	// Set NextPage only if not on last page
	if offset < pagination.TotalPages {
		nextPage := offset + 1
		pagination.NextPage = &nextPage
	}

	// Set PrevPage only if not on first page
	if offset > 1 && (offset - 1) < pagination.TotalPages {
		prevPage := offset - 1
		pagination.PrevPage = &prevPage
	}

	response := &BooksPaginationResponse{
		Books:      books,
		Pagination: pagination,
	}

	helpers.OnRequestSuccess(c, response, "Books retrieved successfully", http.StatusOK)
}

// GetBook godoc
// @Summary      Get a book by ID
// @Description  Retrieves a specific book by its ID
// @Tags         books
// @Produce      json
// @Param        id   path      int  true  "Book ID"
// @Success      200  {object}  helpers.ResponseOk{data=models.Book}
// @Failure      400  {object}  helpers.ErrorResponse "Invalid ID"
// @Failure      404  {object}  helpers.ErrorResponse "Book not found"
// @Router       /api/v1/books/{id} [get]
func GetBook(c *gin.Context) {
	id, ok := helpers.ParseRequestBookId(c)
	if !ok {
		return
	}

	// Try to get from cache first
	book, err := cache.GetBook(id)

	if err != nil {
		// If not in cache, get from database
		bookData, err := repo.GetBookByID(id)

		if err != nil {
			helpers.OnRequestError(c, "Book not found", http.StatusNotFound)
			return
		}

		book = &bookData
	}

	// Store in cache for next time
	if err := cache.SetBook(*book); err != nil {
		log.Printf("Failed to cache book: %v", err)
	}

	helpers.OnRequestSuccess(c, book, "Book retrieved successfully", http.StatusOK)
}
