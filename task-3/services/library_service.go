package services

import (
	"errors"
	"library_management/models"
)

// LibraryManager interface defines the contract for library operations
type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
	AddMember(member models.Member)
	GetMember(memberID int) (*models.Member, error)
}

// Library implements the LibraryManager interface
type Library struct {
	Books   map[int]models.Book
	Members map[int]models.Member
}

// NewLibrary creates a new Library instance
func NewLibrary() *Library {
	return &Library{
		Books:   make(map[int]models.Book),
		Members: make(map[int]models.Member),
	}
}

// AddBook adds a new book to the library
func (l *Library) AddBook(book models.Book) {
	book.Status = "Available"
	l.Books[book.ID] = book
}

// RemoveBook removes a book from the library by its ID
func (l *Library) RemoveBook(bookID int) {
	delete(l.Books, bookID)
}

// BorrowBook allows a member to borrow a book if it is available
func (l *Library) BorrowBook(bookID int, memberID int) error {
	book, exists := l.Books[bookID]
	if !exists {
		return errors.New("book not found")
	}

	if book.Status == "Borrowed" {
		return errors.New("book is already borrowed")
	}

	member, exists := l.Members[memberID]
	if !exists {
		return errors.New("member not found")
	}

	// Update book status
	book.Status = "Borrowed"
	l.Books[bookID] = book

	// Add book to member's borrowed books
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.Members[memberID] = member

	return nil
}

// ReturnBook allows a member to return a borrowed book
func (l *Library) ReturnBook(bookID int, memberID int) error {
	book, exists := l.Books[bookID]
	if !exists {
		return errors.New("book not found")
	}

	member, exists := l.Members[memberID]
	if !exists {
		return errors.New("member not found")
	}

	// Check if member has borrowed this book
	bookIndex := -1
	for i, borrowedBook := range member.BorrowedBooks {
		if borrowedBook.ID == bookID {
			bookIndex = i
			break
		}
	}

	if bookIndex == -1 {
		return errors.New("book not borrowed by this member")
	}

	// Update book status
	book.Status = "Available"
	l.Books[bookID] = book

	// Remove book from member's borrowed books
	member.BorrowedBooks = append(member.BorrowedBooks[:bookIndex], member.BorrowedBooks[bookIndex+1:]...)
	l.Members[memberID] = member

	return nil
}

// ListAvailableBooks lists all available books in the library
func (l *Library) ListAvailableBooks() []models.Book {
	var availableBooks []models.Book
	for _, book := range l.Books {
		if book.Status == "Available" {
			availableBooks = append(availableBooks, book)
		}
	}
	return availableBooks
}

// ListBorrowedBooks lists all books borrowed by a specific member
func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	member, exists := l.Members[memberID]
	if !exists {
		return []models.Book{}
	}
	return member.BorrowedBooks
}

// AddMember adds a new member to the library
func (l *Library) AddMember(member models.Member) {
	member.BorrowedBooks = []models.Book{}
	l.Members[member.ID] = member
}

// GetMember retrieves a member by ID
func (l *Library) GetMember(memberID int) (*models.Member, error) {
	member, exists := l.Members[memberID]
	if !exists {
		return nil, errors.New("member not found")
	}
	return &member, nil
}