package cache

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/aadarshvelu/bms/app/models"
	"github.com/aadarshvelu/bms/config"
)

const (
	BooksCacheKey = "books:all"
	BookKeyPrefix = "book:"
	DefaultTTL    = 3 * time.Hour
)

// SetBooks caches all books
func SetBooks(books []models.Book) error {
	data, err := json.Marshal(books)
	if err != nil {
		return err
	}

	return config.RedisClient.Set(context.Background(), BooksCacheKey, data, DefaultTTL).Err()
}

// GetBooks retrieves paginated books from cache
func GetBooks(limit, offset int) ([]models.Book, int, error) {
	data, err := config.RedisClient.Get(context.Background(), BooksCacheKey).Bytes()
	
	if err != nil {
		return nil, 0, err
	}

	var allBooks []models.Book
	if err = json.Unmarshal(data, &allBooks); err != nil {
		return nil, 0, err
	}

	// Calculate start and end indices for pagination
	startIndex := (offset - 1) * limit
	endIndex := startIndex + limit

	// Validate indices
	if startIndex >= len(allBooks) {
		return []models.Book{}, 0, nil
	}

	if endIndex > len(allBooks) {
		endIndex = len(allBooks)
	}

	return allBooks[startIndex:endIndex], len(allBooks), nil
}

// InvalidateBooks removes the books cache
func InvalidateBooks() error {
	return config.RedisClient.Del(context.Background(), BooksCacheKey).Err()
}

// SetBook caches a single book
func SetBook(book models.Book) error {
	data, err := json.Marshal(book)
	if err != nil {
		return err
	}

	key := BookKeyPrefix + strconv.FormatUint(uint64(book.ID), 10)
	return config.RedisClient.Set(context.Background(), key, data, DefaultTTL).Err()
}

// GetBook retrieves a book from cache by ID
func GetBook(id int) (*models.Book, error) {
	key := BookKeyPrefix + strconv.FormatUint(uint64(id), 10)
	data, err := config.RedisClient.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, err
	}

	var book models.Book
	err = json.Unmarshal(data, &book)
	return &book, err
}

// InvalidateBook removes a single book from cache
func InvalidateBook(id int) error {
	key := BookKeyPrefix + strconv.FormatUint(uint64(id), 10)
	return config.RedisClient.Del(context.Background(), key).Err()
}
