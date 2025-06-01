package agentic

import (
	"context"
	"os"

	"google.golang.org/genai"
)

var (
	Ctx    context.Context
	Client *genai.Client
)

func init() {
	Ctx = context.Background()
	Client, _ = genai.NewClient(Ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
}
