package cli

import (
	"log/slog"
	"os"
)

func init() {
	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, opts)))
}
