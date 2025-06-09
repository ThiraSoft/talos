package talos

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
	// SetLogLevel(slog.LevelError) // Initialize the logger with default level

	Ctx = context.Background()
	Client, _ = genai.NewClient(Ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
}
