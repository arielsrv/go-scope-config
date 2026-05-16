// Package local provides an opt-in autoload variant that, in addition to the
// standard common + scope merge, also merges an optional machine-local override
// file (config.local.yaml / .yml) into the global Viper instance.
//
// Use it via blank import, side-by-side with (or instead of) the regular
// autoload package:
//
//	import _ "github.com/arielsrv/go-scope-config/autoload/local"
//
// Effective precedence (highest wins):
//
//	env vars > config.local > config.[scope] > config.common
//
// The local override file is intended for per-developer settings (localhost
// ports, local DB credentials, mocks, ...) and is typically gitignored. If the
// file does not exist, the autoloader silently skips it.
package local

import (
	"github.com/spf13/viper"

	goscopeconfig "github.com/arielsrv/go-scope-config"
)

func init() { //nolint:gochecknoinits // autoload package is designed to trigger configuration loading via init()
	l := goscopeconfig.NewWithViper(viper.GetViper(), goscopeconfig.WithLocalOverride())
	_ = l.Load()
}
