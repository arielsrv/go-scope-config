// Package goscopeconfig provides a simple and flexible way to load configurations
// based on an environment scope (e.g., dev, prod, staging).
//
// It uses spf13/viper under the hood and allows for configuration inheritance
// by merging a common configuration file with scope-specific overrides.
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
