# go-scope-config

Go package for loading environment-based configurations (`SCOPE`) using [Viper](https://github.com/spf13/viper).

## Installation

```bash
go get github.com/arielsrv/go-scope-config
```

## Usage

The package automatically looks for the `SCOPE` environment variable. If it's not defined, it uses `dev` by default.

It searches for files with the pattern `config.[SCOPE].yaml` or `config.[SCOPE].yml` in a folder (default is `config/`).

### Configuration Merging (Inheritance)

If a `config.common.yaml` (or `.yml`) exists in the configuration directory, it will be loaded first as a base configuration. Then, the scope-specific file (`config.[SCOPE].yaml`) will be merged into it, overriding any shared keys.

### Example File Structure

```text
.
├── config/
│   ├── config.common.yaml (base values)
│   ├── config.dev.yaml    (overrides for dev)
│   └── config.prod.yml    (overrides for prod)
└── main.go
```

### Code

#### Using LoadDefault (Autoloader)

The quickest way to start with default options:

```go
package main

import (
    "fmt"
    "log/slog"
    "os"
    goscopeconfig "github.com/arielsrv/go-scope-config"
)

func main() {
    // Local logger to avoid using the global slog
    logger := slog.Default()

    // Loads from "config/" folder using SCOPE env var (defaults to "dev")
    v, err := goscopeconfig.LoadDefault()
    if err != nil {
        logger.Error("Error loading default config", "error", err)
        os.Exit(1)
    }

    fmt.Printf("App Name: %s\n", v.GetString("app.name"))
}
```

### Automatic Autoloader (init function)

The package provides two ways to automatically load the configuration.

#### 1. Using DefaultViper (named import)

If you import the package with a name, you can access the pre-initialized `DefaultViper` instance.

```go
package main

import (
    "fmt"
    "log/slog"
    "os"
    goscopeconfig "github.com/arielsrv/go-scope-config"
)

func main() {
    // Local logger to avoid using the global slog
    logger := slog.Default()

    // Check if there was an error during automatic loading
    if goscopeconfig.LoadError != nil {
        logger.Error("Error loading config", "error", goscopeconfig.LoadError)
        os.Exit(1)
    }

    // Use the pre-loaded instance
    v := goscopeconfig.DefaultViper
    fmt.Printf("App Name: %s\n", v.GetString("app.name"))
}
```

#### 2. Using Blank Import (global Viper)

If you want to use the standard `github.com/spf13/viper` package directly, you can use the `autoload` sub-package with a blank import. This will load the configuration into Viper's global instance.

```go
package main

import (
    "fmt"
    _ "github.com/arielsrv/go-scope-config/autoload"
    "github.com/spf13/viper"
)

func main() {
    // Configuration is already loaded into global viper
    fmt.Printf("App Name: %s\n", viper.GetString("app.name"))
}
```

#### Manual Initialization (With Options)

```go
package main

import (
    "fmt"
    "log/slog"
    "os"
    goscopeconfig "github.com/arielsrv/go-scope-config"
)

func main() {
    // Local logger to avoid using the global slog
    logger := slog.Default()

    // Initializes the loader. 
    // By default, it looks in the "config/" folder and reads the "SCOPE" environment variable.
    loader := goscopeconfig.New()

    // Optionally, you can configure a custom path
    // loader := goscopeconfig.New(goscopeconfig.WithConfigDir("custom_configs"))

    if err := loader.Load(); err != nil {
        logger.Error("Error loading config", "error", err)
        os.Exit(1)
    }

    // Access values via Viper
    v := loader.Viper()
    fmt.Printf("App Name: %s\n", v.GetString("app.name"))
    fmt.Printf("Current Scope: %s\n", loader.GetScope())
}
```

## Examples

You can find complete examples in the `examples/` folder.

### Run simple example
```bash
go run examples/simple/main.go
```

### Run example with custom directory
```bash
go run examples/custom-dir/main.go
```

### Run example with common config merging
```bash
go run examples/merge-common/main.go
```

### Run autoloader example
```bash
go run examples/autoloader/main.go
```

### Run automatic autoloader example (named import)
```bash
go run examples/automatic/main.go
```

### Run automatic autoloader example (blank import)
```bash
go run examples/blank-import/main.go
```

### Run example with logger (slog JSON)
```bash
go run examples/with-logger/main.go
```

## Logger Support

You can provide a logger that satisfies the `Logger` interface (which has a `Printf(format string, v ...any)` method). This allows easy integration with both the standard `log` package and modern loggers like `slog`.

```go
// Example with slog (JSON format)
handler := slog.NewJSONHandler(os.Stdout, nil)
logger := slog.New(handler)

// Small wrapper to satisfy the Printf interface
type slogWrapper struct { logger *slog.Logger }
func (s *slogWrapper) Printf(f string, v ...any) { s.logger.Info(fmt.Sprintf(f, v...)) }

loader := goscopeconfig.New(
    goscopeconfig.WithLogger(&slogWrapper{logger: logger}),
)
loader.Load()
```

## Environment Variables

In addition to `SCOPE`, the package enables Viper's `AutomaticEnv()`, so you can override any value from the YAML files using environment variables with the corresponding prefix (by default no prefix, with the `.` separator replaced by `_`).

Example: `APP_NAME=my-app` will override the value of `app.name` in the YAML.

## Tests

```bash
go test ./...
```
