package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/spf13/viper"
	"go.uber.org/fx"

	goscopeconfig "github.com/arielsrv/go-scope-config"
)

// NewConfigLoader provides the configuration loader.
func NewConfigLoader() (*viper.Viper, error) {
	loader := goscopeconfig.New(
		goscopeconfig.WithConfigDir("examples/uber-fx/config"),
	)
	err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf("could not load configuration: %w", err)
	}
	return loader.Viper(), nil
}

// NewHTTPServer builds and returns an HTTP server.
func NewHTTPServer(lc fx.Lifecycle, v *viper.Viper) *http.Server {
	port := v.GetInt("app.port")
	if port == 0 {
		port = 8080 // default
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "pong")
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			var listenConfig net.ListenConfig
			ln, err := listenConfig.Listen(ctx, "tcp", server.Addr)
			if err != nil {
				return fmt.Errorf("could not listen on %s: %w", server.Addr, err)
			}
			fmt.Printf("Starting HTTP server at %s\n", server.Addr)
			go func() {
				err = server.Serve(ln)
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					fmt.Printf("HTTP server error: %v\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

	return server
}

func main() {
	app := fx.New(
		fx.Provide(
			NewConfigLoader,
			NewHTTPServer,
		),
		fx.Invoke(func(*http.Server) {}),
	)

	// In a real application, you would use app.Run()
	// But for this example, we'll just start it and then stop it after a quick check
	// or just leave it as is so the user can run it.

	app.Run()
}
