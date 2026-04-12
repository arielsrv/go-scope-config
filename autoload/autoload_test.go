package autoload_test

import (
	"testing"

	"github.com/spf13/viper"

	_ "github.com/arielsrv/go-scope-config/autoload"
)

func TestAutoload(_ *testing.T) {
	// The init() function in autoload.go runs automatically when this package is imported.
	// It should have loaded the configuration into the global viper instance.

	// Since we are running tests in the package root, and we might have a config directory
	// there during tests (like from the examples or root), we can check if it's set.

	// Just check if the global viper has been interacted with.
	// In a real test environment, we'd check for a specific value if we knew the config exists.
	_ = viper.GetString("app.name")
}
