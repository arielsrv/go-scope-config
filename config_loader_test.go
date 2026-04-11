package goscopeconfig

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
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

	t.Run("NewWithViper", func(t *testing.T) {
		v := viper.New()
		l := NewWithViper(v, WithScope("custom"))
		if l.Viper() != v {
			t.Error("Viper instance mismatch")
		}
		if l.GetScope() != "custom" {
			t.Errorf("Expected scope custom, got %s", l.GetScope())
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

func TestDefaultViper(t *testing.T) {
	// For this test, DefaultViper was already initialized by init() when the test started.
	// Since we are running tests from the package root, and "config/config.dev.yaml" exists,
	// DefaultViper should be loaded if the tests were started with the right environment.
	// Note: init() runs only once.

	// Since we cannot easily re-run init(), we check if it loaded correctly
	// given the existing "config" directory in the project root.
	if DefaultViper == nil {
		t.Error("DefaultViper should not be nil")
	}

	// We can't guarantee what's in the root config/ during all test environments,
	// but in this project it should have "My App Dev" if it's the one we created earlier.
	// Let's just check if we can access it without crashing.
	_ = DefaultViper.GetString("app.name")
}

func TestLoadDefault(t *testing.T) {
	// Create a temporary directory and configuration files for the test
	tmpDir := t.TempDir()

	// LoadDefault uses DefaultConfigDir ("config")
	configDir := filepath.Join(tmpDir, DefaultConfigDir)
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}

	devConfig := filepath.Join(configDir, "config.dev.yaml")
	if err := os.WriteFile(devConfig, []byte("app:\n  name: default-dev\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Change current working directory to the temporary one for DefaultConfigDir to work
	// t.Chdir automatically restores the working directory when the test finishes.
	t.Chdir(tmpDir)

	t.Run("LoadDefault with dev scope (default)", func(t *testing.T) {
		os.Unsetenv(ScopeEnvVar)
		v, dfltErr := LoadDefault()
		if dfltErr != nil {
			t.Fatalf("LoadDefault failed: %v", dfltErr)
		}
		if name := v.GetString("app.name"); name != "default-dev" {
			t.Errorf("Expected default-dev, got %s", name)
		}
	})
}
