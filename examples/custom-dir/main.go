package main

import (
	"fmt"
	"log/slog"
	"os"

	goscopeconfig "github.com/arielsrv/go-scope-config"
)

func main() {
	// Local logger to avoid using the global slog
	logger := slog.Default()

	// We force a directory and a specific scope for this example
	loader := goscopeconfig.New(
		goscopeconfig.WithConfigDir("examples/custom-dir/configs"),
		goscopeconfig.WithScope("staging"),
	)

	err := loader.Load()
	if err != nil {
		logger.Error("Error loading config", "error", err)
		os.Exit(1)
	}

	v := loader.Viper()
	fmt.Printf("--- Custom Directory Example ---\n")
	fmt.Printf("Scope: %s\n", loader.GetScope())
	fmt.Printf("DB Host: %s\n", v.GetString("database.host"))
	fmt.Printf("DB User: %s\n", v.GetString("database.user"))
	fmt.Printf("File used: %s\n", loader.GetConfigPath())
}
