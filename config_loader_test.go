package goscopeconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	// Clear environment variables for a clean test
	os.Unsetenv(ScopeEnvVar)

	t.Run("Default ConfigLoader", func(t *testing.T) {
		l := New()
		if l.GetScope() != DefaultScope {
			t.Errorf("Expected scope %s, got %s", DefaultScope, l.GetScope())
		}
	})

	t.Run("With Environmental SCOPE", func(t *testing.T) {
		t.Setenv(ScopeEnvVar, "prod")
		defer os.Unsetenv(ScopeEnvVar)
		l := New()
		if l.GetScope() != "prod" {
			t.Errorf("Expected scope prod, got %s", l.GetScope())
		}
	})

	t.Run("With Custom Option WithScope", func(t *testing.T) {
		l := New(WithScope("staging"))
		if l.GetScope() != "staging" {
			t.Errorf("Expected scope staging, got %s", l.GetScope())
		}
	})
}

func TestLoad(t *testing.T) {
	// Create a temporary directory for tests
	tmpDir := t.TempDir()

	// Create test configuration files
	devConfig := filepath.Join(tmpDir, "config.dev.yaml")
	if err := os.WriteFile(devConfig, []byte("app:\n  name: test-dev\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	prodConfig := filepath.Join(tmpDir, "config.prod.yml")
	if err := os.WriteFile(prodConfig, []byte("app:\n  name: test-prod\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Run("Load dev config", func(t *testing.T) {
		l := New(WithConfigDir(tmpDir), WithScope("dev"))
		if err := l.Load(); err != nil {
			t.Fatalf("Could not load configuration: %v", err)
		}
		if name := l.Viper().GetString("app.name"); name != "test-dev" {
			t.Errorf("Expected test-dev, got %s", name)
		}
	})

	t.Run("Load prod config (yml extension)", func(t *testing.T) {
		l := New(WithConfigDir(tmpDir), WithScope("prod"))
		if err := l.Load(); err != nil {
			t.Fatalf("Could not load configuration: %v", err)
		}
		if name := l.Viper().GetString("app.name"); name != "test-prod" {
			t.Errorf("Expected test-prod, got %s", name)
		}
	})

	t.Run("Config not found error", func(t *testing.T) {
		l := New(WithConfigDir(tmpDir), WithScope("missing"))
		if err := l.Load(); err == nil {
			t.Error("An error was expected when the configuration file does not exist")
		}
	})

	t.Run("Merge with common config", func(t *testing.T) {
		// Create common config file
		commonConfig := filepath.Join(tmpDir, "config.common.yaml")
		commonContent := `
app:
  name: base-app
  version: 1.0.0
database:
  host: localhost
`
		if err := os.WriteFile(commonConfig, []byte(commonContent), 0o644); err != nil {
			t.Fatal(err)
		}

		// Create dev config file that overrides some values
		devConfig = filepath.Join(tmpDir, "config.dev.yaml")
		devContent := `
app:
  name: dev-app
`
		if err := os.WriteFile(devConfig, []byte(devContent), 0o644); err != nil {
			t.Fatal(err)
		}

		l := New(WithConfigDir(tmpDir), WithScope("dev"))
		if err := l.Load(); err != nil {
			t.Fatalf("Could not load configuration: %v", err)
		}

		v := l.Viper()
		// Overridden value
		if name := v.GetString("app.name"); name != "dev-app" {
			t.Errorf("Expected dev-app, got %s", name)
		}
		// Common value (inherited)
		if version := v.GetString("app.version"); version != "1.0.0" {
			t.Errorf("Expected 1.0.0, got %s", version)
		}
		// Another common value
		if host := v.GetString("database.host"); host != "localhost" {
			t.Errorf("Expected localhost, got %s", host)
		}
	})
}
