package main

import (
	"fmt"
	"log"
	"net/http"
	"student-management/config"
	"student-management/routes"
)

func main() {
	// Initialize the database connection
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	
	// Setup routes
	router := routes.SetupRouter(db)
	
	// Start the server
	port := config.GetEnv("PORT", "8080")
	fmt.Printf("Server running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
} 