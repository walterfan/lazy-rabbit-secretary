package book

import (
	"errors"
	"time"

	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"github.com/walterfan/lazy-rabbit-secretary/pkg/database"
	"gorm.io/gorm"
)

// BookRepository provides data access for Book entities
type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository() *BookRepository {
	return &BookRepository{db: database.GetDB()}
}

func (r *BookRepository) Create(book *models.Book) error {
	return r.db.Create(book).Error
}

func (r *BookRepository) GetByID(id string) (*models.Book, error) {
	var book models.Book
	if err := r.db.Preload("Tags").First(&book, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) GetByISBN(isbn string) (*models.Book, error) {
	var book models.Book
	if err := r.db.Preload("Tags").First(&book, "isbn = ?", isbn).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) Update(book *models.Book) error {
	if book.ID == "" {
		return errors.New("missing id for update")
	}
	return r.db.Save(book).Error
}

func (r *BookRepository) Delete(id string) error {
	return r.db.Delete(&models.Book{}, "id = ?", id).Error
}

// Search returns books under a realm with optional filters and pagination
func (r *BookRepository) Search(realmID, query, author, tags string, page, pageSize int) ([]models.Book, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var books []models.Book
	var total int64

	// Build query
	db := r.db.Model(&models.Book{}).Where("realm_id = ?", realmID)

	// Apply filters
	if query != "" {
		db = db.Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%")
	}

	if author != "" {
		db = db.Where("author ILIKE ?", "%"+author+"%")
	}

	if tags != "" {
		// Search by tags using JOIN
		db = db.Joins("JOIN book_tags ON books.id = book_tags.book_id").
			Where("book_tags.name ILIKE ?", "%"+tags+"%").
			Group("books.id")
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results with tags
	if err := db.Preload("Tags").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&books).Error; err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

// GetBorrowedBooks returns books that are currently borrowed
func (r *BookRepository) GetBorrowedBooks(realmID string) ([]models.Book, error) {
	var books []models.Book
	err := r.db.Preload("Tags").
		Where("realm_id = ? AND borrow_time IS NOT NULL AND return_time IS NULL", realmID).
		Order("borrow_time DESC").
		Find(&books).Error
	return books, err
}

// GetOverdueBooks returns books that are past their deadline
func (r *BookRepository) GetOverdueBooks(realmID string) ([]models.Book, error) {
	var books []models.Book
	err := r.db.Preload("Tags").
		Where("realm_id = ? AND borrow_time IS NOT NULL AND return_time IS NULL AND deadline < ?", realmID, "NOW()").
		Order("deadline ASC").
		Find(&books).Error
	return books, err
}

// GetAvailableBooks returns books that are not currently borrowed
func (r *BookRepository) GetAvailableBooks(realmID string) ([]models.Book, error) {
	var books []models.Book
	err := r.db.Preload("Tags").
		Where("realm_id = ? AND (borrow_time IS NULL OR return_time IS NOT NULL)", realmID).
		Order("name ASC").
		Find(&books).Error
	return books, err
}

// BorrowBook marks a book as borrowed
func (r *BookRepository) BorrowBook(bookID string, deadline *time.Time) error {
	return r.db.Model(&models.Book{}).
		Where("id = ?", bookID).
		Updates(map[string]interface{}{
			"borrow_time": time.Now(),
			"return_time": nil,
			"deadline":    deadline,
		}).Error
}

// ReturnBook marks a book as returned
func (r *BookRepository) ReturnBook(bookID string) error {
	return r.db.Model(&models.Book{}).
		Where("id = ?", bookID).
		Updates(map[string]interface{}{
			"return_time": time.Now(),
		}).Error
}

// GetBooksByTag returns books that have a specific tag
func (r *BookRepository) GetBooksByTag(realmID, tagName string) ([]models.Book, error) {
	var books []models.Book
	err := r.db.Preload("Tags").
		Joins("JOIN book_tags ON books.id = book_tags.book_id").
		Where("books.realm_id = ? AND book_tags.name = ?", realmID, tagName).
		Group("books.id").
		Order("books.name ASC").
		Find(&books).Error
	return books, err
}

// GetAllTags returns all unique tags for books in a realm
func (r *BookRepository) GetAllTags(realmID string) ([]string, error) {
	var tags []string
	err := r.db.Model(&models.BookTag{}).
		Joins("JOIN books ON book_tags.book_id = books.id").
		Where("books.realm_id = ?", realmID).
		Distinct("book_tags.name").
		Pluck("book_tags.name", &tags).Error
	return tags, err
}
