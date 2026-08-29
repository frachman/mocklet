package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/frachman/mocklet/apps/api/internal/mocklet"
)

func main() {
	port := getenv("API_PORT", "8080")
	databaseURL := getenv("DATABASE_URL", "postgres://mocklet:mocklet_dev@localhost:5432/mocklet?sslmode=disable")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	repo, err := mocklet.NewPostgresRepository(ctx, databaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer repo.Close()

	server := &http.Server{Addr: ":" + port, Handler: mocklet.NewHandler(repo), ReadHeaderTimeout: 5 * time.Second}
	log.Printf("mocklet api listening on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
