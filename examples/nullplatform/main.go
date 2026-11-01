package main

import (
	"fmt"
	"os"

	goscopeconfig "github.com/arielsrv/go-scope-config"
)

// Parameters that nullplatform injects into every deployed application.
// See https://docs.nullplatform.com/docs/parameters/injected-parameters
//
// When running locally (outside nullplatform) they are not present, so this
// example simulates a deployment to the "production" scope. Inside a real
// nullplatform scope, the injected values take precedence.
var simulatedParameters = map[string]string{ //nolint:gochecknoglobals // example-only simulated environment
	"NP_SCOPE":          "production",
	"NP_APPLICATION":    "users-api",
	"NP_NAMESPACE":      "platform",
	"NP_RELEASE_SEMVER": "1.1.34",
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	err := simulateNullplatform()
	if err != nil {
		return err
	}

	fmt.Println("--- nullplatform Example ---")

	// NP_SCOPE holds the slug of the scope the application is deployed to.
	// The loader reads config.common.yaml and merges config.<NP_SCOPE>.yaml on top.
	// Without NP_SCOPE (e.g. local development) it falls back to "dev".
	loader := goscopeconfig.New(
		goscopeconfig.WithScopeEnv("NP_SCOPE"),
		goscopeconfig.WithConfigDir("nullplatform/config"),
	)

	err = loader.Load()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	v := loader.Viper()
	fmt.Printf("Scope: %s\n", loader.GetScope())
	fmt.Printf("File used: %s\n", loader.GetConfigPath())
	fmt.Printf("DB Host: %s\n", v.GetString("database.host"))
	fmt.Printf("DB Pool: %d\n", v.GetInt("database.pool"))

	// AutomaticEnv maps "np.application" to NP_APPLICATION, so the other
	// injected parameters can be read through Viper as well.
	fmt.Printf("Application: %s\n", v.GetString("np.application"))
	fmt.Printf("Namespace: %s\n", v.GetString("np.namespace"))
	fmt.Printf("Release: %s\n", v.GetString("np.release.semver"))

	return nil
}

// simulateNullplatform sets the injected parameters only when they are missing,
// so the example behaves the same locally and inside a real nullplatform scope.
func simulateNullplatform() error {
	for key, value := range simulatedParameters {
		if _, ok := os.LookupEnv(key); ok {
			continue
		}

		err := os.Setenv(key, value)
		if err != nil {
			return fmt.Errorf("error setting %s: %w", key, err)
		}
	}

	return nil
}
