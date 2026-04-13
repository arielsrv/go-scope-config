package goscopeconfig

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/spf13/viper"
)

// Defaults contain the default configuration values.
type Defaults struct {
	ConfigDir  string
	Scope      string
	ScopeEnv   string
	CommonName string
}

// DefaultConfig returns the default configuration values.
func DefaultConfig() Defaults {
	return Defaults{
		ConfigDir:  "config",
		Scope:      "dev",
		ScopeEnv:   "SCOPE",
		CommonName: "config.common",
	}
}

var defaultConfig = DefaultConfig()

// ConfigLoader handles loading configurations based on the scope.
type ConfigLoader struct {
	logger     Logger
	viper      *viper.Viper
	configDir  string
	scope      string
	scopeEnv   string
	configPath string
}

// Logger is a simple interface for logging.
// Many logging libraries satisfy this interface (e.g., standard [log.Logger], logrus).
type Logger interface {
	Printf(format string, v ...any)
}

// Option defines a function to configure the ConfigLoader.
type Option func(*ConfigLoader)

// WithConfigDir allows specifying a custom folder for the configurations.
func WithConfigDir(dir string) Option {
	return func(loader *ConfigLoader) {
		loader.configDir = dir
	}
}

// WithScope allows forcing a specific scope, ignoring the environment variable.
func WithScope(scope string) Option {
	return func(loader *ConfigLoader) {
		loader.scope = scope
	}
}

// WithScopeEnv allows specifying a custom environment variable to load the scope from.
// Default is SCOPE.
func WithScopeEnv(envName string) Option {
	return func(loader *ConfigLoader) {
		loader.scopeEnv = envName
	}
}

// WithLogger allows providing a logger to show loaded files.
func WithLogger(logger Logger) Option {
	return func(loader *ConfigLoader) {
		loader.logger = logger
	}
}

// New creates a new ConfigLoader instance with the provided options.
func New(opts ...Option) *ConfigLoader {
	return NewWithViper(viper.New(), opts...)
}

// NewWithViper creates a new ConfigLoader instance using a specific Viper instance.
func NewWithViper(viper *viper.Viper, opts ...Option) *ConfigLoader {
	loader := &ConfigLoader{
		viper:     viper,
		configDir: defaultConfig.ConfigDir,
		scopeEnv:  defaultConfig.ScopeEnv,
	}

	for opt := range slices.Values(opts) {
		opt(loader)
	}

	// If the scope was not forced via Option, we look for it in the environment.
	if loader.scope == "" {
		loader.scope = strings.ToLower(cmp.Or(os.Getenv(loader.scopeEnv), defaultConfig.Scope))
	}

	return loader
}

// Load loads the configuration according to the current scope.
// It looks for config.common.yaml/yml first and then merges with config.[scope].yaml/yml.
// The scope-specific configuration overrides values in the common one.
func (r *ConfigLoader) Load() error {
	r.viper.SetConfigType("yaml") // Viper automatically supports YAML and yml if YAML is specified as the type
	r.viper.AddConfigPath(r.configDir)

	// We also allow reading from general environment variables
	r.viper.AutomaticEnv()
	r.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if r.logger != nil {
		r.logger.Printf("Loading configuration for scope: %s (from environment variable: %s)", r.scope, r.scopeEnv)
		r.logger.Printf("Automatic environment variables enabled (replacing '.' with '_')")
		r.logger.Printf("Configuration directory: %s", r.configDir)
	}

	// 1. Try to load config.common.yaml
	r.viper.SetConfigName(defaultConfig.CommonName)
	err := r.viper.ReadInConfig()
	if err != nil {
		// It's okay if the common config doesn't exist
		if _, ok := errors.AsType[viper.ConfigFileNotFoundError](err); !ok {
			return fmt.Errorf("error reading common configuration: %w", err)
		}
	} else {
		if r.logger != nil {
			r.logger.Printf("Loaded common configuration from: %s", r.viper.ConfigFileUsed())
		}
	}

	// 2. Merge with config.[scope].yaml
	scopeConfigName := fmt.Sprintf("config.%s", r.scope)
	r.viper.SetConfigName(scopeConfigName)

	err = r.viper.MergeInConfig()
	if err != nil {
		return fmt.Errorf("error merging configuration file for scope %s in %s: %w", r.scope, r.configDir, err)
	}

	r.configPath = r.viper.ConfigFileUsed()
	if r.logger != nil {
		r.logger.Printf("Loaded scope-specific configuration (%s) from: %s", r.scope, r.configPath)
	}
	return nil
}

var (
	// DefaultViper is the automatically loaded Viper instance.
	// It is initialized during package loading via init().
	DefaultViper *viper.Viper
	// ErrLoad stores any error encountered during automatic loading.
	ErrLoad error
)

func init() {
	DefaultViper, ErrLoad = LoadDefault()
}

// LoadDefault is a shortcut that creates a new ConfigLoader with default options,
// loads the configuration, and returns the internal Viper instance.
// This is useful for quick starts where default behavior is enough.
func LoadDefault() (*viper.Viper, error) {
	loader := New()
	err := loader.Load()
	if err != nil {
		return nil, err
	}
	return loader.Viper(), nil
}

// Viper returns the internal viper instance to access values.
func (r *ConfigLoader) Viper() *viper.Viper {
	return r.viper
}

// GetScope returns the current scope being used.
func (r *ConfigLoader) GetScope() string {
	return r.scope
}

// GetConfigPath returns the absolute path of the loaded configuration file (the scope-specific one).
func (r *ConfigLoader) GetConfigPath() string {
	return r.configPath
}
