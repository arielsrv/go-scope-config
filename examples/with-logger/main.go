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

func (r *slogWrapper) Printf(format string, v ...any) {
	r.logger.Info(fmt.Sprintf(format, v...))
}

func main() {
	fmt.Println("--- Example with slog (JSON format) ---")

	// Create a new slog logger with JSON handler
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Wrap it to satisfy the Logger interface
	loader := goscopeconfig.New(
		goscopeconfig.WithScope("dev"),
		goscopeconfig.WithLogger(&slogWrapper{logger: logger}),
	)

	err := loader.Load()
	if err != nil {
		logger.Error("Error loading config", "error", err)
		os.Exit(1)
	}

	v := loader.Viper()
	fmt.Printf("App Name: %s\n", v.GetString("app.name"))
	fmt.Printf("Port: %s\n", v.GetString("app.port"))
}
