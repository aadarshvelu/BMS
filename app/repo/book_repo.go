package repo

import (
	"github.com/aadarshvelu/bms/app/models"
	"github.com/aadarshvelu/bms/config"
)

type CreateBookPayload struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
	Year   *int    `json:"year" binding:"required"`
}

type UpdateBookPayload struct {
	Title  *string `json:"title,omitempty"`
	Author *string `json:"author,omitempty"`
	Year   *int    `json:"year,omitempty"`
}

func CreateBook(bookReq CreateBookPayload) (models.Book, error) {
	book := models.Book{
		Title:  bookReq.Title,
		Author: bookReq.Author,
		Year:   *bookReq.Year,
	}

	result := config.DB.Create(&book)

	return book, result.Error
}

func UpdateBook(id int, bookReq UpdateBookPayload) (models.Book, error) {
	var book models.Book

	if err := config.DB.First(&book, id).Error; err != nil {
		return book, err
	}

	// Only update fields that are provided in the request
	if bookReq.Title != nil {
		book.Title = *bookReq.Title
	}

	if bookReq.Author != nil {
		book.Author = *bookReq.Author
	}

	if bookReq.Year != nil {
		book.Year = *bookReq.Year
	}

	result := config.DB.Save(&book)

	return book, result.Error
}

func  DeleteBook(id int) (int64, error) {
	result := config.DB.Delete(&models.Book{}, id)
	return result.RowsAffected, result.Error
}

func GetAllBooks() ([]models.Book, error) {
	var books []models.Book

	result := config.DB.Find(&books)

	return books, result.Error
}

func GetBookByID(id int) (models.Book, error) {
	var book models.Book

	result := config.DB.First(&book, id)

	return book, result.Error
}