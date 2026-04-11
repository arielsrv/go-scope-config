package main

import (
	"fmt"
	"log/slog"
	"os"

	goscopeconfig "github.com/arielsrv/go-scope-config"
)

// slogWrapper adapts slog to the goscopeconfig.Logger interface.
type slogWrapper struct {
	logger *slog.Logger
}

func (s *slogWrapper) Printf(format string, v ...any) {
	s.logger.Info(fmt.Sprintf(format, v...))
}

func main() {
	fmt.Println("--- Example with slog (JSON format) ---")

	// Create a new slog logger with JSON handler
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)

	// Wrap it to satisfy the Logger interface
	loader := goscopeconfig.New(
		goscopeconfig.WithScope("dev"),
		goscopeconfig.WithLogger(&slogWrapper{logger: logger}),
	)

	if err := loader.Load(); err != nil {
		logger.Error("Error loading config", "error", err)
		os.Exit(1)
	}

	v := loader.Viper()
	fmt.Printf("App Name: %s\n", v.GetString("app.name"))
	fmt.Printf("Port: %s\n", v.GetString("app.port"))
}
