# Library Management System Documentation

## Overview
This is a console-based library management system built in Go that demonstrates the use of structs, interfaces, and other Go functionalities such as methods, slices, and maps.

## Architecture

### Project Structure
```
library_management/
├── main.go                      # Entry point of the application
├── controllers/
│   └── library_controller.go   # Handles console input and service invocation
├── models/
│   ├── book.go                 # Book struct definition
│   └── member.go               # Member struct definition
├── services/
│   └── library_service.go     # Business logic and data manipulation
├── docs/
│   └── documentation.md        # System documentation
└── go.mod                      # Module definition
```

## Core Components

### Models

#### Book Struct
```go
type Book struct {
    ID     int
    Title  string
    Author string
    Status string // "Available" or "Borrowed"
}
```

#### Member Struct
```go
type Member struct {
    ID            int
    Name          string
    BorrowedBooks []Book
}
```

### Interfaces

#### LibraryManager Interface
Defines the contract for library operations:
- `AddBook(book Book)`
- `RemoveBook(bookID int)`
- `BorrowBook(bookID int, memberID int) error`
- `ReturnBook(bookID int, memberID int) error`
- `ListAvailableBooks() []Book`
- `ListBorrowedBooks(memberID int) []Book`

### Services

#### Library Struct
Implements the LibraryManager interface with:
- `Books map[int]Book` - stores all books with ID as key
- `Members map[int]Member` - stores all members with ID as key

## Features

### Book Management
- **Add Book**: Add new books to the library inventory
- **Remove Book**: Remove books from the library by ID
- **List Available Books**: Display all books currently available for borrowing

### Member Management
- **Add Member**: Register new library members
- **List Borrowed Books**: View books borrowed by a specific member

### Borrowing System
- **Borrow Book**: Allow members to borrow available books
- **Return Book**: Process book returns and update availability

## Error Handling

The system includes comprehensive error handling for:
- Book not found scenarios
- Member not found scenarios
- Attempting to borrow already borrowed books
- Attempting to return books not borrowed by the member

## Usage

### Running the Application
```bash
cd task-3
go run main.go
```

### Sample Data
The system comes pre-loaded with sample books and members:

**Books:**
- The Go Programming Language by Alan Donovan
- Clean Code by Robert Martin
- Design Patterns by Gang of Four

**Members:**
- John Doe (ID: 1)
- Jane Smith (ID: 2)

### Console Interface
The application provides an interactive menu with the following options:
1. Add a new book
2. Remove a book
3. Borrow a book
4. Return a book
5. List available books
6. List borrowed books by member
7. Add a new member
8. Exit

## Technical Implementation

### Data Storage
- Uses Go maps for efficient O(1) lookup operations
- Books stored with integer ID as key
- Members stored with integer ID as key
- Borrowed books maintained as slices within member structs

### Interface Implementation
- Clean separation of concerns through interface design
- Service layer implements business logic
- Controller layer handles user interaction
- Model layer defines data structures

### Memory Management
- Efficient use of Go's built-in data structures
- Proper slice manipulation for borrowed books tracking
- No external dependencies required

## Future Enhancements

Potential improvements could include:
- Persistent data storage (database integration)
- Book search functionality
- Due date tracking for borrowed books
- Fine calculation system
- Member borrowing limits
- Book reservation system