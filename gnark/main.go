package main

import (
	"fmt"
	"log"
	"net/http"
	"ppba_project/gnark/api"
	"ppba_project/gnark/db"
)

func main() {
	// Initialize the database
	db, err := database.NewDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	// Setup the HTTP server
	http.HandleFunc("/enroll", api.HandleEnroll(db))
	http.HandleFunc("/verify", api.HandleVerify(db))
	http.HandleFunc("/retrieve", api.HandleRetrieve(db))

	port := ":8080"
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
