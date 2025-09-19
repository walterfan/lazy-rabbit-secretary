package book

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"gorm.io/gorm"
)

// BookService contains business logic for books
type BookService struct {
	repo *BookRepository
}

func NewBookService(repo *BookRepository) *BookService {
	return &BookService{
		repo: repo,
	}
}

// CreateBookRequest defines the allowed input for creating a book
type CreateBookRequest struct {
	ISBN        string    `json:"isbn" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Tags        string    `json:"tags"` // comma-separated tags
	Deadline    time.Time `json:"deadline"`
}

// UpdateBookRequest defines the allowed input for updating a book
type UpdateBookRequest struct {
	ISBN        string    `json:"isbn"`
	Name        string    `json:"name"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Tags        string    `json:"tags"` // comma-separated tags
	Deadline    time.Time `json:"deadline"`
}

// BookResponse defines the output format for book data
type BookResponse struct {
	ID          string    `json:"id"`
	RealmID     string    `json:"realm_id"`
	ISBN        string    `json:"isbn"`
	Name        string    `json:"name"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	BorrowTime  time.Time `json:"borrow_time"`
	ReturnTime  time.Time `json:"return_time"`
	Deadline    time.Time `json:"deadline"`
	Tags        []string  `json:"tags"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedBy   string    `json:"updated_by"`
	UpdatedTime time.Time `json:"updated_time"`
}

// BookListResponse defines the response format for book lists
type BookListResponse struct {
	Books []BookResponse `json:"books"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

// CreateFromInput creates a new book from input data
func (s *BookService) CreateFromInput(req *CreateBookRequest, realmID, createdBy string) (*BookResponse, error) {
	// Validate request
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if ISBN already exists
	existingBook, err := s.repo.GetByISBN(req.ISBN)
	if err == nil && existingBook != nil {
		return nil, fmt.Errorf("book with ISBN %s already exists", req.ISBN)
	}

	// Create book model
	book := &models.Book{
		ID:          uuid.New().String(),
		RealmID:     realmID,
		ISBN:        req.ISBN,
		Name:        req.Name,
		Author:      req.Author,
		Description: req.Description,
		Price:       req.Price,
		Deadline:    &req.Deadline,
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
	}

	// Parse and create tags
	if req.Tags != "" {
		tagNames := s.parseTags(req.Tags)
		book.Tags = s.createBookTags(book.ID, tagNames, createdBy)
	}

	// Save to database
	if err := s.repo.Create(book); err != nil {
		return nil, fmt.Errorf("failed to create book: %w", err)
	}

	return s.toResponse(book), nil
}

// UpdateFromInput updates a book from input data
func (s *BookService) UpdateFromInput(id string, req *UpdateBookRequest, updatedBy string) (*BookResponse, error) {
	// Validate request
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get existing book
	book, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("book not found")
		}
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	// Update fields if provided
	if req.ISBN != "" {
		// Check if new ISBN already exists (excluding current book)
		existingBook, err := s.repo.GetByISBN(req.ISBN)
		if err == nil && existingBook != nil && existingBook.ID != id {
			return nil, fmt.Errorf("book with ISBN %s already exists", req.ISBN)
		}
		book.ISBN = req.ISBN
	}
	if req.Name != "" {
		book.Name = req.Name
	}
	if req.Author != "" {
		book.Author = req.Author
	}
	if req.Description != "" {
		book.Description = req.Description
	}
	if req.Price >= 0 {
		book.Price = req.Price
	}
	if !req.Deadline.IsZero() {
		book.Deadline = &req.Deadline
	}

	// Update tags if provided
	if req.Tags != "" {
		tagNames := s.parseTags(req.Tags)
		book.Tags = s.createBookTags(book.ID, tagNames, updatedBy)
	}

	book.UpdatedBy = updatedBy

	// Save to database
	if err := s.repo.Update(book); err != nil {
		return nil, fmt.Errorf("failed to update book: %w", err)
	}

	return s.toResponse(book), nil
}

// SearchBooks searches for books with filters and pagination
func (s *BookService) SearchBooks(realmID, query, author, tags string, page, pageSize int) (*BookListResponse, error) {
	books, total, err := s.repo.Search(realmID, query, author, tags, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to search books: %w", err)
	}

	responses := make([]BookResponse, len(books))
	for i, book := range books {
		responses[i] = *s.toResponse(&book)
	}

	return &BookListResponse{
		Books: responses,
		Total: total,
		Page:  page,
		Limit: pageSize,
	}, nil
}

// GetBookByID retrieves a book by ID
func (s *BookService) GetBookByID(id string) (*BookResponse, error) {
	book, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("book not found")
		}
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	return s.toResponse(book), nil
}

// DeleteBook deletes a book by ID
func (s *BookService) DeleteBook(id string) error {
	// Check if book exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("book not found")
		}
		return fmt.Errorf("failed to get book: %w", err)
	}

	// Check if book is currently borrowed
	book, _ := s.repo.GetByID(id)
	if book.BorrowTime != nil && book.ReturnTime == nil {
		return fmt.Errorf("cannot delete borrowed book")
	}

	return s.repo.Delete(id)
}

// BorrowBook marks a book as borrowed
func (s *BookService) BorrowBook(id string, deadline time.Time) (*BookResponse, error) {
	book, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("book not found")
		}
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	// Check if book is already borrowed
	if book.BorrowTime != nil && book.ReturnTime == nil {
		return nil, fmt.Errorf("book is already borrowed")
	}

	// Borrow the book
	if err := s.repo.BorrowBook(id, &deadline); err != nil {
		return nil, fmt.Errorf("failed to borrow book: %w", err)
	}

	// Get updated book
	updatedBook, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated book: %w", err)
	}

	return s.toResponse(updatedBook), nil
}

// ReturnBook marks a book as returned
func (s *BookService) ReturnBook(id string) (*BookResponse, error) {
	book, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("book not found")
		}
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	// Check if book is borrowed
	if book.BorrowTime == nil || book.ReturnTime != nil {
		return nil, fmt.Errorf("book is not currently borrowed")
	}

	// Return the book
	if err := s.repo.ReturnBook(id); err != nil {
		return nil, fmt.Errorf("failed to return book: %w", err)
	}

	// Get updated book
	updatedBook, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated book: %w", err)
	}

	return s.toResponse(updatedBook), nil
}

// GetBorrowedBooks returns all currently borrowed books
func (s *BookService) GetBorrowedBooks(realmID string) ([]BookResponse, error) {
	books, err := s.repo.GetBorrowedBooks(realmID)
	if err != nil {
		return nil, fmt.Errorf("failed to get borrowed books: %w", err)
	}

	responses := make([]BookResponse, len(books))
	for i, book := range books {
		responses[i] = *s.toResponse(&book)
	}

	return responses, nil
}

// GetOverdueBooks returns all overdue books
func (s *BookService) GetOverdueBooks(realmID string) ([]BookResponse, error) {
	books, err := s.repo.GetOverdueBooks(realmID)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue books: %w", err)
	}

	responses := make([]BookResponse, len(books))
	for i, book := range books {
		responses[i] = *s.toResponse(&book)
	}

	return responses, nil
}

// GetAvailableBooks returns all available books
func (s *BookService) GetAvailableBooks(realmID string) ([]BookResponse, error) {
	books, err := s.repo.GetAvailableBooks(realmID)
	if err != nil {
		return nil, fmt.Errorf("failed to get available books: %w", err)
	}

	responses := make([]BookResponse, len(books))
	for i, book := range books {
		responses[i] = *s.toResponse(&book)
	}

	return responses, nil
}

// GetBooksByTag returns books with a specific tag
func (s *BookService) GetBooksByTag(realmID, tagName string) ([]BookResponse, error) {
	books, err := s.repo.GetBooksByTag(realmID, tagName)
	if err != nil {
		return nil, fmt.Errorf("failed to get books by tag: %w", err)
	}

	responses := make([]BookResponse, len(books))
	for i, book := range books {
		responses[i] = *s.toResponse(&book)
	}

	return responses, nil
}

// GetAllTags returns all unique tags in a realm
func (s *BookService) GetAllTags(realmID string) ([]string, error) {
	return s.repo.GetAllTags(realmID)
}

// Helper methods

func (s *BookService) validateCreateRequest(req *CreateBookRequest) error {
	if req.ISBN == "" {
		return errors.New("ISBN is required")
	}
	if req.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

func (s *BookService) validateUpdateRequest(req *UpdateBookRequest) error {
	// All fields are optional for updates
	return nil
}

func (s *BookService) parseTags(tags string) []string {
	if tags == "" {
		return []string{}
	}

	tagList := strings.Split(tags, ",")
	result := make([]string, 0, len(tagList))
	for _, tag := range tagList {
		trimmed := strings.TrimSpace(tag)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func (s *BookService) createBookTags(bookID string, tagNames []string, createdBy string) []models.BookTag {
	tags := make([]models.BookTag, len(tagNames))
	for i, name := range tagNames {
		tags[i] = models.BookTag{
			ID:        uuid.New().String(),
			BookID:    bookID,
			Name:      name,
			CreatedBy: createdBy,
			UpdatedBy: createdBy,
		}
	}
	return tags
}

func (s *BookService) toResponse(book *models.Book) *BookResponse {
	tags := make([]string, len(book.Tags))
	for i, tag := range book.Tags {
		tags[i] = tag.Name
	}

	response := &BookResponse{
		ID:          book.ID,
		RealmID:     book.RealmID,
		ISBN:        book.ISBN,
		Name:        book.Name,
		Author:      book.Author,
		Description: book.Description,
		Price:       book.Price,
		Tags:        tags,
		CreatedBy:   book.CreatedBy,
		CreatedAt:   book.CreatedAt,
		UpdatedBy:   book.UpdatedBy,
		UpdatedTime: book.UpdatedTime,
	}

	if book.BorrowTime != nil {
		response.BorrowTime = *book.BorrowTime
	}
	if book.ReturnTime != nil {
		response.ReturnTime = *book.ReturnTime
	}
	if book.Deadline != nil {
		response.Deadline = *book.Deadline
	}

	return response
}
