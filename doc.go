// Package goscopeconfig provides a simple and flexible way to load configurations
// based on an environment scope (e.g., dev, prod, staging).
//
// It uses spf13/viper under the hood and allows for configuration inheritance
// by merging a common configuration file with scope-specific overrides.
//
// Optionally, a machine-local override file (config.local.yaml / .yml) can be
// merged on top via WithLocalOverride(). This file is intended for per-developer
// settings (localhost ports, local DB credentials, mocks, ...) and is typically
// gitignored. If it does not exist, Load() does not fail.
//
// Effective precedence (highest wins):
//
//	env vars > config.local > config.[scope] > config.common
//
// Basic usage:
//
//	loader := goscopeconfig.New()
//	if err := loader.Load(); err != nil {
//	    log.Fatal(err)
//	}
//	v := loader.Viper()
//	fmt.Println(v.GetString("app.name"))
//
// With machine-local overrides:
//
//	loader := goscopeconfig.New(goscopeconfig.WithLocalOverride())
//	_ = loader.Load()
//
// Using the autoloader:
//
//	v, err := goscopeconfig.LoadDefault()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(v.GetString("app.name"))
//
// For automatic initialization using the global Viper instance, see the autoload subpackage.
package goscopeconfig
