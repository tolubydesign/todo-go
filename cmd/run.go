/*
Copyright © 2025 Tolu Adesina <tolu.adesina...>
*/
package cmd

import (
	"context"
	"fmt"
	"net"

	"net/http"

	"github.com/spf13/cobra"
	configuration "github.com/tolubydesign/todo-go/app/config"
	h "github.com/tolubydesign/todo-go/app/handler"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Handler h.Handler
type EchoHandler h.EchoHandler

// ServerParams defines the dependencies for the HTTP server setup.
type ServerParams struct {
	fx.In
	Lifecycle fx.Lifecycle   // The fx lifecycle hook
	Mux       *http.ServeMux // The configured router
}

// NewHTTPServer builds an HTTP server that will begin serving requests
// when the Fx application starts.
func NewHTTPServer(p ServerParams, mux *http.ServeMux, log *zap.Logger) {
	c, _ := configuration.GetConfiguration()
	if c.Configuration.Port == "" {
		log.Fatal("Port value not specified")
	}

	var port = c.Configuration.Port
	log.Info("Server boot on port", zap.String("port", port))
	var addr string = fmt.Sprintf(":%s", port)

	server := &http.Server{
		Addr:    addr,
		Handler: p.Mux,
		BaseContext: func(net.Listener) context.Context {
			return context.Background()
		},
	}

	// Add server startup/shutdown hooks to the fx.Lifecycle
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Info("Starting HTTP server on.")
			ln, err := net.Listen("tcp", server.Addr)
			if err != nil {
				log.Fatal("HTTP server failed: ", zap.String("error", err.Error()))
				return err
			}
			log.Info("Starting HTTP server at", zap.String("addr", server.Addr))
			go server.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// This runs when the fx.App stops. On SIGINT.
			log.Warn("Stopping HTTP server...")
			return server.Shutdown(ctx)
		},
	})
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the ToDo Backend API",
	Long:  "This ToDo application uses Cobra for CLI commands and uber-fx for dependency injection and lifecycle management.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")

		_, err := configuration.BuildConfiguration()
		if err != nil {
			fmt.Println("error", err.Error())
			panic(err)
		}

		app := fx.New(
			// Provide the dependencies
			fx.Provide(
				h.HandlerInstance,
				h.ProvideMux,
				h.NewEchoHandler,
				zap.NewExample,
			),

			// invoke background tasks
			fx.Invoke(NewHTTPServer),
			// TODO: Invoke init mysql database
		)

		// starting the application
		app.Run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
