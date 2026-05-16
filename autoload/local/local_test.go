package local_test

import (
	"testing"

	"github.com/spf13/viper"

	_ "github.com/arielsrv/go-scope-config/autoload/local"
)

func TestAutoloadLocal(t *testing.T) {
	t.Parallel()
	// The init() function in the local autoload package runs automatically when
	// this package is imported (blank import above). It loads common + scope +
	// optional config.local into the global viper instance.
	//
	// Just verify the global viper is usable and didn't panic during init.
	_ = viper.GetString("app.name")
}
