package controllers

import (
	"bufio"
	"fmt"
	"library_management/models"
	"library_management/services"
	"os"
	"strconv"
	"strings"
)

// LibraryController handles console input and invokes service methods
type LibraryController struct {
	libraryService services.LibraryManager
	scanner        *bufio.Scanner
}

// NewLibraryController creates a new LibraryController instance
func NewLibraryController(libraryService services.LibraryManager) *LibraryController {
	return &LibraryController{
		libraryService: libraryService,
		scanner:        bufio.NewScanner(os.Stdin),
	}
}

// Start begins the console interface
func (lc *LibraryController) Start() {
	fmt.Println("Welcome to the Library Management System!")
	fmt.Println("=========================================")

	// Add some sample data
	lc.addSampleData()

	for {
		lc.showMenu()
		choice := lc.getInput("Enter your choice: ")

		switch choice {
		case "1":
			lc.addBook()
		case "2":
			lc.removeBook()
		case "3":
			lc.borrowBook()
		case "4":
			lc.returnBook()
		case "5":
			lc.listAvailableBooks()
		case "6":
			lc.listBorrowedBooks()
		case "7":
			lc.addMember()
		case "8":
			fmt.Println("Thank you for using the Library Management System!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
		fmt.Println()
	}
}

func (lc *LibraryController) showMenu() {
	fmt.Println("\n--- Library Management System ---")
	fmt.Println("1. Add a new book")
	fmt.Println("2. Remove a book")
	fmt.Println("3. Borrow a book")
	fmt.Println("4. Return a book")
	fmt.Println("5. List available books")
	fmt.Println("6. List borrowed books by member")
	fmt.Println("7. Add a new member")
	fmt.Println("8. Exit")
}

func (lc *LibraryController) getInput(prompt string) string {
	fmt.Print(prompt)
	lc.scanner.Scan()
	return strings.TrimSpace(lc.scanner.Text())
}

func (lc *LibraryController) getIntInput(prompt string) (int, error) {
	input := lc.getInput(prompt)
	return strconv.Atoi(input)
}

func (lc *LibraryController) addBook() {
	fmt.Println("\n--- Add New Book ---")
	
	id, err := lc.getIntInput("Enter book ID: ")
	if err != nil {
		fmt.Println("Invalid ID. Please enter a number.")
		return
	}

	title := lc.getInput("Enter book title: ")
	author := lc.getInput("Enter book author: ")

	book := models.Book{
		ID:     id,
		Title:  title,
		Author: author,
	}

	lc.libraryService.AddBook(book)
	fmt.Printf("Book '%s' by %s has been added successfully!\n", title, author)
}

func (lc *LibraryController) removeBook() {
	fmt.Println("\n--- Remove Book ---")
	
	id, err := lc.getIntInput("Enter book ID to remove: ")
	if err != nil {
		fmt.Println("Invalid ID. Please enter a number.")
		return
	}

	lc.libraryService.RemoveBook(id)
	fmt.Printf("Book with ID %d has been removed successfully!\n", id)
}

func (lc *LibraryController) borrowBook() {
	fmt.Println("\n--- Borrow Book ---")
	
	bookID, err := lc.getIntInput("Enter book ID to borrow: ")
	if err != nil {
		fmt.Println("Invalid book ID. Please enter a number.")
		return
	}

	memberID, err := lc.getIntInput("Enter member ID: ")
	if err != nil {
		fmt.Println("Invalid member ID. Please enter a number.")
		return
	}

	err = lc.libraryService.BorrowBook(bookID, memberID)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	fmt.Printf("Book with ID %d has been borrowed successfully!\n", bookID)
}

func (lc *LibraryController) returnBook() {
	fmt.Println("\n--- Return Book ---")
	
	bookID, err := lc.getIntInput("Enter book ID to return: ")
	if err != nil {
		fmt.Println("Invalid book ID. Please enter a number.")
		return
	}

	memberID, err := lc.getIntInput("Enter member ID: ")
	if err != nil {
		fmt.Println("Invalid member ID. Please enter a number.")
		return
	}

	err = lc.libraryService.ReturnBook(bookID, memberID)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	fmt.Printf("Book with ID %d has been returned successfully!\n", bookID)
}

func (lc *LibraryController) listAvailableBooks() {
	fmt.Println("\n--- Available Books ---")
	
	books := lc.libraryService.ListAvailableBooks()
	if len(books) == 0 {
		fmt.Println("No books are currently available.")
		return
	}

	fmt.Printf("%-5s %-30s %-20s %-10s\n", "ID", "Title", "Author", "Status")
	fmt.Println(strings.Repeat("-", 70))
	for _, book := range books {
		fmt.Printf("%-5d %-30s %-20s %-10s\n", book.ID, book.Title, book.Author, book.Status)
	}
}

func (lc *LibraryController) listBorrowedBooks() {
	fmt.Println("\n--- Borrowed Books by Member ---")
	
	memberID, err := lc.getIntInput("Enter member ID: ")
	if err != nil {
		fmt.Println("Invalid member ID. Please enter a number.")
		return
	}

	books := lc.libraryService.ListBorrowedBooks(memberID)
	if len(books) == 0 {
		fmt.Printf("No books are currently borrowed by member %d.\n", memberID)
		return
	}

	fmt.Printf("Books borrowed by member %d:\n", memberID)
	fmt.Printf("%-5s %-30s %-20s\n", "ID", "Title", "Author")
	fmt.Println(strings.Repeat("-", 60))
	for _, book := range books {
		fmt.Printf("%-5d %-30s %-20s\n", book.ID, book.Title, book.Author)
	}
}

func (lc *LibraryController) addMember() {
	fmt.Println("\n--- Add New Member ---")
	
	id, err := lc.getIntInput("Enter member ID: ")
	if err != nil {
		fmt.Println("Invalid ID. Please enter a number.")
		return
	}

	name := lc.getInput("Enter member name: ")

	member := models.Member{
		ID:   id,
		Name: name,
	}

	lc.libraryService.AddMember(member)
	fmt.Printf("Member '%s' has been added successfully!\n", name)
}

func (lc *LibraryController) addSampleData() {
	// Add sample books
	sampleBooks := []models.Book{
		{ID: 1, Title: "The Go Programming Language", Author: "Alan Donovan"},
		{ID: 2, Title: "Clean Code", Author: "Robert Martin"},
		{ID: 3, Title: "Design Patterns", Author: "Gang of Four"},
	}

	for _, book := range sampleBooks {
		lc.libraryService.AddBook(book)
	}

	// Add sample members
	sampleMembers := []models.Member{
		{ID: 1, Name: "John Doe"},
		{ID: 2, Name: "Jane Smith"},
	}

	for _, member := range sampleMembers {
		lc.libraryService.AddMember(member)
	}
}