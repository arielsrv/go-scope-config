package goscopeconfig

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	DefaultConfigDir = "config"
	DefaultScope     = "dev"
	ScopeEnvVar      = "SCOPE"
	CommonConfigName = "config.common"
)

// ConfigLoader handles loading configurations based on the scope.
type ConfigLoader struct {
	v          *viper.Viper
	configDir  string
	scope      string
	configPath string
}

// Option defines a function to configure the ConfigLoader.
type Option func(*ConfigLoader)

// WithConfigDir allows specifying a custom folder for the configurations.
func WithConfigDir(dir string) Option {
	return func(l *ConfigLoader) {
		l.configDir = dir
	}
}

// WithScope allows forcing a specific scope, ignoring the environment variable.
func WithScope(scope string) Option {
	return func(l *ConfigLoader) {
		l.scope = scope
	}
}

// New creates a new ConfigLoader instance with the provided options.
func New(opts ...Option) *ConfigLoader {
	l := &ConfigLoader{
		v:         viper.New(),
		configDir: DefaultConfigDir,
	}

	for _, opt := range opts {
		opt(l)
	}

	// If the scope was not forced via Option, we look for it in the environment.
	if l.scope == "" {
		l.scope = strings.ToLower(os.Getenv(ScopeEnvVar))
		if l.scope == "" {
			l.scope = DefaultScope
		}
	}

	return l
}

// Load loads the configuration according to the current scope.
// It looks for config.common.yaml/yml first and then merges with config.[scope].yaml/yml.
// The scope-specific configuration overrides values in the common one.
func (l *ConfigLoader) Load() error {
	l.v.SetConfigType("yaml") // Viper automatically supports yaml and yml if yaml is specified as the type
	l.v.AddConfigPath(l.configDir)

	// We also allow reading from general environment variables
	l.v.AutomaticEnv()
	l.v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 1. Try to load config.common.yaml
	l.v.SetConfigName(CommonConfigName)
	if err := l.v.ReadInConfig(); err != nil {
		// It's okay if the common config doesn't exist
		if _, ok := errors.AsType[viper.ConfigFileNotFoundError](err); !ok {
			return fmt.Errorf("error reading common configuration: %w", err)
		}
	}

	// 2. Merge with config.[scope].yaml
	scopeConfigName := fmt.Sprintf("config.%s", l.scope)
	l.v.SetConfigName(scopeConfigName)

	if err := l.v.MergeInConfig(); err != nil {
		return fmt.Errorf("error merging configuration file for scope %s in %s: %w", l.scope, l.configDir, err)
	}

	l.configPath = l.v.ConfigFileUsed()
	return nil
}

// Viper returns the internal viper instance to access values.
func (l *ConfigLoader) Viper() *viper.Viper {
	return l.v
}

// GetScope returns the current scope being used.
func (l *ConfigLoader) GetScope() string {
	return l.scope
}

// GetConfigPath returns the absolute path of the loaded configuration file (the scope-specific one).
func (l *ConfigLoader) GetConfigPath() string {
	return l.configPath
}
