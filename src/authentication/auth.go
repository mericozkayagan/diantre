package authentication

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/v50/github"
	"github.com/joho/godotenv"
)

func Auth() (*github.Client, context.Context) {
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}
	ctx := context.Background()
	client := github.NewTokenClient(ctx, token)

	return client, ctx
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
