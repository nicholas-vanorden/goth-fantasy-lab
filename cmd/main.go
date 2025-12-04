package main

import (
	"context"
	"fmt"
	"goth-ffb-players/internal/handlers"
	"goth-ffb-players/internal/services/auth"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	clientID := os.Getenv("OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("OAUTH_CLIENT_SECRET")
	serverPort := os.Getenv("SERVER_PORT")
	redirectUrl := fmt.Sprintf("https://%s:%s/oauth/callback", os.Getenv("SERVER"), serverPort)
	authUrl := os.Getenv("OAUTH_AUTH_URL")
	tokenUrl := os.Getenv("OAUTH_TOKEN_URL")
	certdotpem := os.Getenv("TLS_CERT_PATH")
	keydotpem := os.Getenv("TLS_KEY_PATH")

	authCache := auth.New(clientID, clientSecret, redirectUrl, authUrl, tokenUrl, "token.json")
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Authentication
	mux.HandleFunc("GET /oauth/login", handlers.OAuthLogin(authCache))
	mux.HandleFunc("GET /oauth/callback", handlers.OAuthCallback(authCache))
	mux.HandleFunc("GET /players", handlers.Players(authCache))
	mux.HandleFunc("GET /players/{id}", handlers.PlayerDetail(authCache))
	mux.HandleFunc("GET /search", handlers.SearchPlayers(authCache))
	mux.Handle("/", http.RedirectHandler("/players", http.StatusSeeOther))

	server := &http.Server{
		Addr:    ":" + serverPort,
		Handler: mux,
	}

	// Create a channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start HTTPS server in a goroutine
	go func() {
		log.Println("Server starting on ", server.Addr)
		err := server.ListenAndServeTLS(certdotpem, keydotpem)
		if err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Block until a signal is received
	<-quit
	log.Println("Shutting down server...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("Server gracefully stopped.")
}
