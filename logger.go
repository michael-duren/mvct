package mvct

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
)

type LoggerConfig struct {
	Path     string
	LogLevel *log.Level
}

// InitLogger initializes the global Logger to write to specified file
// To use it use the 'slog' logger
func (a *Application[M]) UseLogger(config LoggerConfig) error {
	logDir, level := config.Path, config.LogLevel
	// Determine log file path
	if logDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get user home dir: %w", err)
		}
		logDir = filepath.Join(wd, "logs")
	}

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log dir: %w", err)
	}

	logFile := filepath.Join(logDir, fmt.Sprintf("g7c-%s.log", time.Now().Format("2006-01-02")))
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	charmHandler := log.New(io.Writer(f))
	charmHandler.SetReportCaller(true)

	if level != nil {
		charmHandler.SetLevel(*level)
	} else {
		charmHandler.SetLevel(log.InfoLevel)
	}

	logger := slog.New(charmHandler)
	slog.SetDefault(logger)

	slog.Info("Logger initialized", "log_file", logFile)
	return nil
}

func LogLevel(level log.Level) *log.Level {
	return &level
}
