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

	// SCOPE=dev by default if not defined
	loader := goscopeconfig.New()

	err := loader.Load()
	if err != nil {
		logger.Error("Error loading config", "error", err)
		os.Exit(1)
	}

	v := loader.Viper()
	fmt.Printf("--- Simple Example ---\n")
	fmt.Printf("Scope: %s\n", loader.GetScope())
	fmt.Printf("App Name: %s\n", v.GetString("app.name"))
	fmt.Printf("Port: %d\n", v.GetInt("app.port"))
	fmt.Printf("File used: %s\n", loader.GetConfigPath())
}
