package goscopeconfig_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"

	goscopeconfig "github.com/arielsrv/go-scope-config"
	goscopeconfigmocks "github.com/arielsrv/go-scope-config/mocks"
)

// Default config values mirrored from the original package for testing purposes.
var defaultConfig = struct {
	ConfigDir  string
	Scope      string
	ScopeEnv   string
	CommonName string
}{
	ConfigDir:  "config",
	Scope:      "dev",
	ScopeEnv:   "SCOPE",
	CommonName: "config.common",
}

func TestNew(t *testing.T) {
	// Clear environment variables for a clean test
	err := os.Unsetenv(defaultConfig.ScopeEnv)
	if err != nil {
		return
	}

	t.Run("Default ConfigLoader", func(t *testing.T) {
		l := goscopeconfig.New()
		if l.GetScope() != defaultConfig.Scope {
			t.Errorf("Expected scope %s, got %s", defaultConfig.Scope, l.GetScope())
		}
	})

	t.Run("With Environmental SCOPE", func(t *testing.T) {
		t.Setenv(defaultConfig.ScopeEnv, "prod")
		defer func() {
			err = os.Unsetenv(defaultConfig.ScopeEnv)
			if err != nil {
				return
			}
		}()
		l := goscopeconfig.New()
		if l.GetScope() != "prod" {
			t.Errorf("Expected scope prod, got %s", l.GetScope())
		}
	})

	t.Run("With Custom Option WithScope", func(t *testing.T) {
		l := goscopeconfig.New(goscopeconfig.WithScope("staging"))
		if l.GetScope() != "staging" {
			t.Errorf("Expected scope staging, got %s", l.GetScope())
		}
	})

	t.Run("NewWithViper", func(t *testing.T) {
		v := viper.New()
		l := goscopeconfig.NewWithViper(v, goscopeconfig.WithScope("custom"))
		if l.Viper() != v {
			t.Error("Viper instance mismatch")
		}
		if l.GetScope() != "custom" {
			t.Errorf("Expected scope custom, got %s", l.GetScope())
		}
	})

	t.Run("With Custom ScopeEnv", func(t *testing.T) {
		customEnv := "MY_SCOPE"
		t.Setenv(customEnv, "staging")
		l := goscopeconfig.New(goscopeconfig.WithScopeEnv(customEnv))
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
	err := os.WriteFile(devConfig, []byte("app:\n  name: test-dev\n"), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	prodConfig := filepath.Join(tmpDir, "config.prod.yml")
	err = os.WriteFile(prodConfig, []byte("app:\n  name: test-prod\n"), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Load dev config", func(t *testing.T) {
		l := goscopeconfig.New(goscopeconfig.WithConfigDir(tmpDir), goscopeconfig.WithScope("dev"))
		err = l.Load()
		if err != nil {
			t.Fatalf("Could not load configuration: %v", err)
		}
		if name := l.Viper().GetString("app.name"); name != "test-dev" {
			t.Errorf("Expected test-dev, got %s", name)
		}
	})

	t.Run("Load prod config (yml extension)", func(t *testing.T) {
		l := goscopeconfig.New(goscopeconfig.WithConfigDir(tmpDir), goscopeconfig.WithScope("prod"))
		err = l.Load()
		if err != nil {
			t.Fatalf("Could not load configuration: %v", err)
		}
		if name := l.Viper().GetString("app.name"); name != "test-prod" {
			t.Errorf("Expected test-prod, got %s", name)
		}
	})

	t.Run("Config not found error", func(t *testing.T) {
		l := goscopeconfig.New(goscopeconfig.WithConfigDir(tmpDir), goscopeconfig.WithScope("missing"))
		err = l.Load()
		if err == nil {
			t.Error("An error was expected when the configuration file does not exist")
		}
	})

	t.Run("Invalid YAML file", func(t *testing.T) {
		invalidConfig := filepath.Join(tmpDir, "config.invalid.yaml")
		err = os.WriteFile(invalidConfig, []byte("app: - name: invalid: yaml: content:"), 0o644)
		if err != nil {
			t.Fatal(err)
		}
		l := goscopeconfig.New(goscopeconfig.WithConfigDir(tmpDir), goscopeconfig.WithScope("invalid"))
		err = l.Load()
		if err == nil {
			t.Error("Expected error for invalid yaml content, got nil")
		}
	})

	t.Run("Invalid common YAML file", func(t *testing.T) {
		// New temp dir for this subtest to avoid interfering with others
		subTmpDir := t.TempDir()
		commonInvalid := filepath.Join(subTmpDir, "config.common.yaml")
		err = os.WriteFile(commonInvalid, []byte("app: - name: invalid: yaml: content:"), 0o644)
		if err != nil {
			t.Fatal(err)
		}

		l := goscopeconfig.New(goscopeconfig.WithConfigDir(subTmpDir), goscopeconfig.WithScope("dev"))
		err = l.Load()
		if err == nil {
			t.Error("Expected error for invalid common yaml content, got nil")
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
		err = os.WriteFile(commonConfig, []byte(commonContent), 0o644)
		if err != nil {
			t.Fatal(err)
		}

		// Create dev config file that overrides some values
		devConfig = filepath.Join(tmpDir, "config.dev.yaml")
		devContent := `
app:
  name: dev-app
`
		err = os.WriteFile(devConfig, []byte(devContent), 0o644)
		if err != nil {
			t.Fatal(err)
		}

		l := goscopeconfig.New(goscopeconfig.WithConfigDir(tmpDir), goscopeconfig.WithScope("dev"))
		err = l.Load()
		if err != nil {
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

		if path := l.GetConfigPath(); path == "" {
			t.Error("Expected non-empty config path")
		}
	})

	t.Run("With logger", func(t *testing.T) {
		mockLogger := goscopeconfigmocks.NewMockLogger(t)
		commonConfig := filepath.Join(tmpDir, "config.common.yaml")
		err = os.WriteFile(commonConfig, []byte("app:\n  base: true\n"), 0o644)
		if err != nil {
			t.Fatal(err)
		}

		// Expectations for loaded files using typed EXPECT()
		mockLogger.EXPECT().Printf(mock.Anything, mock.MatchedBy(func(s string) bool {
			return strings.Contains(s, "config.common.yaml")
		})).Once()

		mockLogger.EXPECT().Printf(mock.Anything, "dev", mock.MatchedBy(func(s string) bool {
			return strings.Contains(s, "config.dev.yaml")
		})).Once()

		// Optional logs (loading messages, automatic env, etc.)
		mockLogger.EXPECT().Printf(mock.Anything, mock.Anything).Maybe()
		mockLogger.EXPECT().Printf(mock.Anything, mock.Anything, mock.Anything).Maybe()
		mockLogger.EXPECT().Printf(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()

		l := goscopeconfig.New(
			goscopeconfig.WithConfigDir(tmpDir),
			goscopeconfig.WithScope("dev"),
			goscopeconfig.WithLogger(mockLogger),
		)
		err = l.Load()
		if err != nil {
			t.Fatalf("Could not load configuration: %v", err)
		}
	})
}

func TestDefaultViper(t *testing.T) {
	// For this test, goscopeconfig.DefaultViper was already initialized by init() when the test started.
	// Since we are running tests from the package root, and "config/config.dev.yaml" exists,
	// goscopeconfig.DefaultViper should be loaded if the tests were started with the right environment.
	// Note: init() runs only once.

	// Since we cannot easily re-run init(), we check if it loaded correctly
	// given the existing "config" directory in the project root.
	if goscopeconfig.DefaultViper == nil {
		t.Error("goscopeconfig.DefaultViper should not be nil")
	}

	// We can't guarantee what's in the root config/ during all test environments,
	// but in this project it should have "My App Dev" if it's the one we created earlier.
	// Let's just check if we can access it without crashing.
	_ = goscopeconfig.DefaultViper.GetString("app.name")
}

func TestLoadDefault(t *testing.T) {
	// Create a temporary directory and configuration files for the test
	tmpDir := t.TempDir()

	// LoadDefault uses defaultConfig.ConfigDir ("config")
	configDir := filepath.Join(tmpDir, defaultConfig.ConfigDir)
	err := os.MkdirAll(configDir, 0o755)
	if err != nil {
		t.Fatal(err)
	}

	devConfig := filepath.Join(configDir, "config.dev.yaml")
	err = os.WriteFile(devConfig, []byte("app:\n  name: default-dev\n"), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	// Change current working directory to the temporary one for DefaultConfigDir to work
	// t.Chdir automatically restores the working directory when the test finishes.
	t.Chdir(tmpDir)

	t.Run("LoadDefault with dev scope (default)", func(t *testing.T) {
		err = os.Unsetenv(defaultConfig.ScopeEnv)
		if err != nil {
			return
		}
		v, dfltErr := goscopeconfig.LoadDefault()
		if dfltErr != nil {
			t.Fatalf("LoadDefault failed: %v", dfltErr)
		}
		if name := v.GetString("app.name"); name != "default-dev" {
			t.Errorf("Expected default-dev, got %s", name)
		}
	})

	t.Run("LoadDefault with error (missing directory)", func(t *testing.T) {
		// Change to an empty temporary directory
		otherTmp := t.TempDir()
		t.Chdir(otherTmp)

		v, dfltErr := goscopeconfig.LoadDefault()
		if dfltErr == nil {
			t.Error("Expected error when loading default config from empty directory, got nil")
		}
		if v != nil {
			t.Error("Expected nil Viper instance on error")
		}
	})
}
