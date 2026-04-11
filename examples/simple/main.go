package main

import (
	"fmt"
	"log"

	goscopeconfig "github.com/arielsrv/go-scope-config"
)

func main() {
	// SCOPE=dev by default if not defined
	loader := goscopeconfig.New()

	if err := loader.Load(); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	v := loader.Viper()
	fmt.Printf("--- Simple Example ---\n")
	fmt.Printf("Scope: %s\n", loader.GetScope())
	fmt.Printf("App Name: %s\n", v.GetString("app.name"))
	fmt.Printf("Port: %d\n", v.GetInt("app.port"))
	fmt.Printf("File used: %s\n", loader.GetConfigPath())
}
