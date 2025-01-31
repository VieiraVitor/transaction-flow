package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func InitLogger() {
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo, // Pode ser ajustado para Debug, Warning, etc.
		AddSource: true,           // Adiciona arquivo e linha do log automaticamente
	}))
}
