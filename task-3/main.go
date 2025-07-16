package main

import (
	"library_management/controllers"
	"library_management/services"
)

func main() {
	// Initialize the library service
	libraryService := services.NewLibrary()

	// Initialize the controller with the service
	controller := controllers.NewLibraryController(libraryService)

	// Start the console interface
	controller.Start()
}