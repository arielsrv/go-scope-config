package autoload

import (
	"github.com/spf13/viper"

	goscopeconfig "github.com/arielsrv/go-scope-config"
)

func init() { //nolint:gochecknoinits // autoload package is designed to trigger configuration loading via init()
	// Load configuration into the global Viper instance
	l := goscopeconfig.NewWithViper(viper.GetViper())
	_ = l.Load()
}
