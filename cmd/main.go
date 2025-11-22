package main

import (
	"log"
	"net/http"

	"goth-ffb-players/internal/handlers"
)

func main() {
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("GET /", handlers.PlayerList)
	mux.HandleFunc("GET /players/{id}", handlers.PlayerDetail)
	mux.HandleFunc("GET /search", handlers.SearchPlayers)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
