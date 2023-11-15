package main

import (
	"egit90/urlShortner/database"
	"egit90/urlShortner/handlers"
	"log"
	"net/http"
)

func main() {
	db, err := database.NewDBManager("./data")
	if err != nil {
		log.Fatalf("Failed to create DBManager: %v", err)
	}
	defer db.CloseDB()

	router := handlers.DefaultMux(db)
	handler := handlers.DbHandler(db, router)

	log.Println("Starting the server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
