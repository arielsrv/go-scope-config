package goscopeconfig

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

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

// ConfigLoader handles loading configurations based on the scope.
type ConfigLoader struct {
	logger     Logger
	viper      *viper.Viper
	configDir  string
	scope      string
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

// WithLogger allows providing a logger to show loaded files.
func WithLogger(logger Logger) Option {
	return func(l *ConfigLoader) {
		l.logger = logger
	}
}

// New creates a new ConfigLoader instance with the provided options.
func New(opts ...Option) *ConfigLoader {
	return NewWithViper(viper.New(), opts...)
}

// NewWithViper creates a new ConfigLoader instance using a specific Viper instance.
func NewWithViper(viper *viper.Viper, opts ...Option) *ConfigLoader {
	l := &ConfigLoader{
		viper:     viper,
		configDir: defaultConfig.ConfigDir,
	}

	for _, opt := range opts {
		opt(l)
	}

	// If the scope was not forced via Option, we look for it in the environment.
	if l.scope == "" {
		l.scope = strings.ToLower(os.Getenv(defaultConfig.ScopeEnv))
		if l.scope == "" {
			l.scope = defaultConfig.Scope
		}
	}

	return l
}

// Load loads the configuration according to the current scope.
// It looks for config.common.yaml/yml first and then merges with config.[scope].yaml/yml.
// The scope-specific configuration overrides values in the common one.
func (r *ConfigLoader) Load() error {
	r.viper.SetConfigType("yaml") // Viper automatically supports yaml and yml if yaml is specified as the type
	r.viper.AddConfigPath(r.configDir)

	// We also allow reading from general environment variables
	r.viper.AutomaticEnv()
	r.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 1. Try to load config.common.yaml
	r.viper.SetConfigName(defaultConfig.CommonName)
	if err := r.viper.ReadInConfig(); err != nil {
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

	if err := r.viper.MergeInConfig(); err != nil {
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
	l := New()
	if err := l.Load(); err != nil {
		return nil, err
	}
	return l.Viper(), nil
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
