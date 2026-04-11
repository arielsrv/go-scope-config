package autoload

import (
	"github.com/spf13/viper"

	goscopeconfig "github.com/arielsrv/go-scope-config"
)

func init() {
	// Load configuration into the global Viper instance
	l := goscopeconfig.NewWithViper(viper.GetViper())
	_ = l.Load()
}
