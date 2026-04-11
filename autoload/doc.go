// Package autoload provides automatic configuration loading into the global Viper instance.
//
// By simply importing this package with a blank import, it will initialize the
// configuration using the default options (SCOPE from environment or "dev", and "config/" directory)
// and load it directly into the viper.GetViper() instance.
//
// Usage:
//
//	import _ "github.com/arielsrv/go-scope-config/autoload"
//	import "github.com/spf13/viper"
//
//	func main() {
//	    fmt.Println(viper.GetString("app.name"))
//	}
package autoload
