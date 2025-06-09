package talos

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func init() {
	SetLogLevel(slog.LevelError) // Définit le niveau de log par défaut à Error
}

// Init initialise un logger configuré
func SetLogLevel(level slog.Leveler) {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	logger = slog.New(handler)
}
