// Example: blank import that ALSO enables the optional config.local override.
//
// Importing autoload/local with a blank import triggers an init() that loads
// common + scope + the optional config.local file into the global Viper instance.
//
// Effective precedence (highest wins):
//
//	env vars > config.local > config.[scope] > config.common
package main

import (
	"fmt"

	"github.com/spf13/viper"

	// This single blank import is enough — it autoloads common + scope + local.
	_ "github.com/arielsrv/go-scope-config/autoload/local"
)

func main() {
	fmt.Println("--- Blank Import + Local Override Example ---")
	fmt.Printf("app.name       = %s\n", viper.GetString("app.name"))
	fmt.Printf("app.port       = %d\n", viper.GetInt("app.port"))
	fmt.Printf("database.host  = %s   (from config.local)\n", viper.GetString("database.host"))
	fmt.Printf("database.port  = %d   (from config.local)\n", viper.GetInt("database.port"))
}
