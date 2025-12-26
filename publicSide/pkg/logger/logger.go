// Package logger предоставляет утилиты для настройки логирования в приложении.
package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Setup настраивает и возвращает новый экземпляр *slog.Logger.
// Он определяет уровень логирования на основе предоставленной строки `logLevel`.
// Поддерживаемые уровни: "DEBUG", "WARN", "ERROR". По умолчанию используется "INFO".
// Логгер выводит текстовые сообщения в os.Stdout.
func Setup(logLevel string) *slog.Logger {
	var level slog.Level
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		level = slog.LevelDebug
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
}
