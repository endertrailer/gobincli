package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"gopbincli/internal/handler"
	"gopbincli/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	ctx := context.Background()
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is empty or not set")
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Warning %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Database ping failed: %v\n", err)
	}
	repo := repository.NewPostgresRepo(pool)

	authHandler := &handler.AuthHandler{Repo: repo}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /create", authHandler.CreateBinHandler)
	mux.HandleFunc("GET /bin/{id}", authHandler.GetBinByIdHandler)
	port := os.Getenv("port")
	if port == "" {
		port = "8080" // Default fallback
	}
	http.ListenAndServe(":"+port, mux)
}
