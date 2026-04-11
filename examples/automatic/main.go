package main

import (
	"fmt"
	"log"

	// Importing the package initializes DefaultViper via init().
	goscopeconfig "github.com/arielsrv/go-scope-config"
)

func main() {
	fmt.Println("--- Automatic Autoloader Example ---")

	// Check if there was an error during automatic loading
	if goscopeconfig.ErrLoad != nil {
		log.Fatalf("Error loading automatic config: %v", goscopeconfig.ErrLoad)
	}

	// Access the pre-loaded Viper instance
	v := goscopeconfig.DefaultViper

	fmt.Printf("App Name: %s\n", v.GetString("app.name"))
	fmt.Printf("App Port: %d\n", v.GetInt("app.port"))
}
