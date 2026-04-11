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

```go
package main

import (
    "fmt"
    "log"
    goscopeconfig "github.com/arielsrv/go-scope-config"
)

func main() {
    // Initializes the loader. 
    // By default, it looks in the "config/" folder and reads the "SCOPE" environment variable.
    loader := goscopeconfig.New()

    // Optionally, you can configure a custom path
    // loader := goscopeconfig.New(goscopeconfig.WithConfigDir("custom_configs"))

    if err := loader.Load(); err != nil {
        log.Fatalf("Error loading config: %v", err)
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

## Environment Variables

In addition to `SCOPE`, the package enables Viper's `AutomaticEnv()`, so you can override any value from the YAML files using environment variables with the corresponding prefix (by default no prefix, with the `.` separator replaced by `_`).

Example: `APP_NAME=my-app` will override the value of `app.name` in the YAML.

## Tests

```bash
go test ./...
```
