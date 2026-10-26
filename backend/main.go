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

	cardHandler := &handlers.CardHandler{
		Store:    dbStore,
		Validate: validate,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /cards", cardHandler.CreateCard)
	mux.HandleFunc("GET /cards", cardHandler.ListCards)
	mux.HandleFunc("GET /cards/{id}", cardHandler.GetCard)
	mux.HandleFunc("PUT /cards/{id}", cardHandler.UpdateCard)
	mux.HandleFunc("POST /cards/{id}/move", cardHandler.MoveCard)
	mux.HandleFunc("DELETE /cards/{id}", cardHandler.DeleteCard)

	columnHandler := &handlers.ColumnHandler{Store: dbStore, Validate: validate}

	mux.HandleFunc("POST /columns", columnHandler.CreateColumn)
	mux.HandleFunc("GET /columns", columnHandler.ListColumns)
	mux.HandleFunc("GET /columns/{id}", columnHandler.GetColumn)
	mux.HandleFunc("PUT /columns/{id}", columnHandler.UpdateColumn)
	mux.HandleFunc("DELETE /columns/{id}", columnHandler.DeleteColumn)

	boardHandler := &handlers.BoardHandler{Store: dbStore, Validate: validate}

	mux.HandleFunc("POST /boards", boardHandler.CreateBoard)
	mux.HandleFunc("GET /boards", boardHandler.ListBoards)
	mux.HandleFunc("GET /boards/{id}", boardHandler.GetBoard)
	mux.HandleFunc("PUT /boards/{id}", boardHandler.UpdateBoard)
	mux.HandleFunc("DELETE /boards/{id}", boardHandler.DeleteBoard)

	userHandler := &handlers.UserHandler{Store: dbStore, Validate: validate}

	mux.HandleFunc("POST /users", userHandler.CreateUser)
	mux.HandleFunc("GET /users", userHandler.ListUsers)
	mux.HandleFunc("GET /users/{id}", userHandler.GetUser)
	mux.HandleFunc("PUT /users/{id}", userHandler.UpdateUser)
	mux.HandleFunc("DELETE /users/{id}", userHandler.DeleteUser)

	authHandler := &handlers.AuthHandler{Store: dbStore, Validate: validate}

	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("POST /logout", authHandler.Logout)
	mux.Handle("GET /me", authHandler.RequireAuth(http.HandlerFunc(authHandler.Me)))

	fmt.Println("Server running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
