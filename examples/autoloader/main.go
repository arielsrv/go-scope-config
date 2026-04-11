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

	// For this example to work with LoadDefault, we need the "config" folder
	// in the root of where we run the program.
	// Since we are running from the project root, it will use the existing "config" folder.

	fmt.Println("--- Autoloader Example ---")

	// LoadDefault uses SCOPE from env or "dev" by default, and "config/" directory.
	v, err := goscopeconfig.LoadDefault()
	if err != nil {
		logger.Error("Error loading default config", "error", err)
		os.Exit(1)
	}

	fmt.Printf("App Name: %s\n", v.GetString("app.name"))
	fmt.Printf("App Port: %d\n", v.GetInt("app.port"))

	scope := os.Getenv("SCOPE")
	if scope == "" {
		scope = "dev (default)"
	}
	fmt.Printf("Loaded scope: %s\n", scope)
}
