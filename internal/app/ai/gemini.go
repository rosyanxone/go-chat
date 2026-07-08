package ai

import (
	"context"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func Client() *genai.Client {
	ctx := context.Background()

	apiKey := os.Getenv("GEMINI_API_KEY")

	if apiKey == "" {
		log.Fatal("Error: GEMINI_API_KEY not exists on the environment variable.")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return client
}
