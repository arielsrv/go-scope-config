package main

import (
	"fmt"
	"os"

	goscopeconfig "github.com/arielsrv/go-scope-config"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// 1. Set a custom environment variable for the scope
	os.Setenv("APP_ENV", "prod")
	defer os.Unsetenv("APP_ENV")

	fmt.Println("--- Custom Scope Environment Variable Example ---")

	// 2. Configure the loader to use "APP_ENV" instead of the default "SCOPE"
	loader := goscopeconfig.New(
		goscopeconfig.WithScopeEnv("APP_ENV"),
		goscopeconfig.WithConfigDir("config"), // Pointing to the root config directory
	)

	// 3. Load the configuration
	err := loader.Load()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	// 4. Verify the loaded scope and values
	fmt.Printf("Detected Scope: %s\n", loader.GetScope())
	fmt.Printf("App Port: %d (should be 80 for prod)\n", loader.Viper().GetInt("app.port"))

	return nil
}
