package main

import (
	"context"
	"darulabror/config"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
)

func main() {
	ctx := context.Background()
	db := config.ConnectionDb()
	validate := validator.New()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
