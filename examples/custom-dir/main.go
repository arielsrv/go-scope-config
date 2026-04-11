package main

import (
	"fmt"
	"log"

	goscopeconfig "github.com/arielsrv/go-scope-config"
)

func main() {
	// We force a directory and a specific scope for this example
	loader := goscopeconfig.New(
		goscopeconfig.WithConfigDir("examples/custom-dir/configs"),
		goscopeconfig.WithScope("staging"),
	)

	if err := loader.Load(); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	v := loader.Viper()
	fmt.Printf("--- Custom Directory Example ---\n")
	fmt.Printf("Scope: %s\n", loader.GetScope())
	fmt.Printf("DB Host: %s\n", v.GetString("database.host"))
	fmt.Printf("DB User: %s\n", v.GetString("database.user"))
	fmt.Printf("File used: %s\n", loader.GetConfigPath())
}
