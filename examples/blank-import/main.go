package main

import (
	"fmt"

	// Importing this sub-package automatically loads config into global viper instance.
	"github.com/spf13/viper"

	_ "github.com/arielsrv/go-scope-config/autoload"
)

func main() {
	fmt.Println("--- Blank Import Autoloader Example ---")

	// We use the standard viper package directly
	fmt.Printf("App Name: %s\n", viper.GetString("app.name"))
	fmt.Printf("App Port: %d\n", viper.GetInt("app.port"))
	fmt.Printf("REST Connections: %d\n", viper.GetInt("rest.connections"))
}
