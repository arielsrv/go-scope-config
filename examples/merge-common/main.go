package main

import (
	"fmt"
	"log"
	"os"

	goscopeconfig "github.com/arielsrv/go-scope-config"
)

func main() {
	// Set SCOPE to prod to see merging and overriding in action
	os.Setenv("SCOPE", "prod")

	// Initialize the loader pointing to the examples/merge-common/configs folder
	loader := goscopeconfig.New(goscopeconfig.WithConfigDir("examples/merge-common/configs"))

	if err := loader.Load(); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	v := loader.Viper()

	fmt.Printf("--- Configuration with Merge and Override ---\n")
	fmt.Printf("Scope: %s\n", loader.GetScope())
	fmt.Printf("Config File Used: %s\n", loader.GetConfigPath())
	fmt.Printf("\n")

	// Values from config.common.yaml
	fmt.Printf("App Version (from common): %s\n", v.GetString("app.version"))
	fmt.Printf("Database Host (from common): %s\n", v.GetString("database.host"))

	// Value from config.common.yaml but overridden by config.prod.yaml
	fmt.Printf("App Name (overridden by prod): %s\n", v.GetString("app.name"))

	// Value only present in config.prod.yaml
	fmt.Printf("Database Port (only in prod): %d\n", v.GetInt("database.port"))
}
