package main

import (
	"fmt"
	"log"
	"os"
	"time"

	goscopeconfig "github.com/arielsrv/go-scope-config"
)

func main() {
	// Load configuration using defaults.
	// It will look for the SCOPE environment variable.
	loader := goscopeconfig.New(
		goscopeconfig.WithLogger(log.Default()), // Use standard logger to see output
	)

	fmt.Println("--- Docker Compose Example ---")

	err := loader.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	v := loader.Viper()

	// Print some values from the configuration
	fmt.Printf("App Name: %s\n", v.GetString("app.name"))
	fmt.Printf("Environment: %s\n", loader.GetScope())
	fmt.Printf("Database Host: %s\n", v.GetString("db.host"))
	fmt.Printf("API Key: %s\n", v.GetString("api.key"))

	fmt.Println("Service is running... (Ctrl+C to stop)")

	// Keep the container alive
	for {
		time.Sleep(10 * time.Second)
	}
}
