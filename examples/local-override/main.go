// Example: local override file (`config.local.yaml`) for per-developer settings.
//
// Layering (highest precedence wins):
//
//	env vars > config.local.yaml > config.[scope].yaml > config.common.yaml
//
// The local file is optional and intended to be gitignored so each developer
// can override values (DB host, ports, feature flags, ...) without touching
// the shared scope configurations.
package main

import (
	"fmt"
	"log"
	"os"

	goscopeconfig "github.com/arielsrv/go-scope-config"
)

func main() {
	logger := log.New(os.Stdout, "[config] ", log.LstdFlags)

	loader := goscopeconfig.New(
		goscopeconfig.WithConfigDir("local-override/config"),
		goscopeconfig.WithLocalOverride(), // opt-in: merge config.local.{yaml,yml} on top
		goscopeconfig.WithLogger(logger),
	)
	if err := loader.Load(); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	v := loader.Viper()
	fmt.Println("---")
	fmt.Printf("scope          = %s\n", loader.GetScope())
	fmt.Printf("app.name       = %s   (from scope dev)\n", v.GetString("app.name"))
	fmt.Printf("app.version    = %s     (inherited from common)\n", v.GetString("app.version"))
	fmt.Printf("database.host  = %s    (overridden by config.local)\n", v.GetString("database.host"))
	fmt.Printf("database.port  = %d      (overridden by config.local)\n", v.GetInt("database.port"))
	fmt.Printf("database.name  = %s     (inherited from common)\n", v.GetString("database.name"))
	fmt.Printf("features.cache = %t       (overridden by config.local)\n", v.GetBool("features.cache"))
}
