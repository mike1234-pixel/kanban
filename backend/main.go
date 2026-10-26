package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"kanban/internal/handlers"
	"kanban/internal/store"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	dbStore, err := store.NewPostgresStore(connStr)
	validate := validator.New()

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	mux := http.NewServeMux()

	authHandler := &handlers.AuthHandler{Store: dbStore, Validate: validate}
	// protected wraps a handler in RequireAuth, so only logged-in users reach it.
	protected := func(h http.HandlerFunc) http.Handler { return authHandler.RequireAuth(h) }

	boardHandler := &handlers.BoardHandler{Store: dbStore, Validate: validate}

	mux.Handle("POST /boards", protected(boardHandler.CreateBoard))
	mux.Handle("GET /boards", protected(boardHandler.ListBoards))
	mux.Handle("GET /boards/{id}", protected(boardHandler.GetBoard))
	mux.Handle("PUT /boards/{id}", protected(boardHandler.UpdateBoard))
	mux.Handle("DELETE /boards/{id}", protected(boardHandler.DeleteBoard))

	columnHandler := &handlers.ColumnHandler{Store: dbStore, Validate: validate}

	mux.Handle("POST /columns", protected(columnHandler.CreateColumn))
	mux.Handle("GET /columns", protected(columnHandler.ListColumns))
	mux.Handle("GET /columns/{id}", protected(columnHandler.GetColumn))
	mux.Handle("PUT /columns/{id}", protected(columnHandler.UpdateColumn))
	mux.Handle("DELETE /columns/{id}", protected(columnHandler.DeleteColumn))

	cardHandler := &handlers.CardHandler{Store: dbStore, Validate: validate}

	mux.Handle("POST /cards", protected(cardHandler.CreateCard))
	mux.Handle("GET /cards", protected(cardHandler.ListCards))
	mux.Handle("GET /cards/{id}", protected(cardHandler.GetCard))
	mux.Handle("PUT /cards/{id}", protected(cardHandler.UpdateCard))
	mux.Handle("POST /cards/{id}/move", protected(cardHandler.MoveCard))
	mux.Handle("DELETE /cards/{id}", protected(cardHandler.DeleteCard))

	userHandler := &handlers.UserHandler{Store: dbStore, Validate: validate}

	mux.HandleFunc("POST /users", userHandler.CreateUser)
	mux.HandleFunc("GET /users", userHandler.ListUsers)
	mux.HandleFunc("GET /users/{id}", userHandler.GetUser)
	mux.HandleFunc("PUT /users/{id}", userHandler.UpdateUser)
	mux.HandleFunc("DELETE /users/{id}", userHandler.DeleteUser)

	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("POST /logout", authHandler.Logout)
	mux.Handle("GET /me", protected(authHandler.Me))

	fmt.Println("Server running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
