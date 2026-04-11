package main

import (
	"fmt"
	"log/slog"
	"os"

	// Importing the package initializes DefaultViper via init().
	goscopeconfig "github.com/arielsrv/go-scope-config"
)

func main() {
	// Local logger to avoid using the global slog
	logger := slog.Default()

	fmt.Println("--- Automatic Autoloader Example ---")

	// Check if there was an error during automatic loading
	if goscopeconfig.ErrLoad != nil {
		logger.Error("Error loading automatic config", "error", goscopeconfig.ErrLoad)
		os.Exit(1)
	}

	// Access the pre-loaded Viper instance
	v := goscopeconfig.DefaultViper

	fmt.Printf("App Name: %s\n", v.GetString("app.name"))
	fmt.Printf("App Port: %d\n", v.GetInt("app.port"))
}
