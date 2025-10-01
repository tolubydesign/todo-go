/*
Copyright © 2025 Tolu Adesina <tolu.adesina...>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"net"

	"net/http"

	"github.com/spf13/cobra"
	configuration "github.com/tolubydesign/todo-go/app/config"
	"github.com/tolubydesign/todo-go/app/db"
	h "github.com/tolubydesign/todo-go/app/handler"
	"github.com/tolubydesign/todo-go/app/logging"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ServerParams defines the dependencies for the HTTP server setup.
type ServerParams struct {
	fx.In
	Lifecycle fx.Lifecycle // The fx lifecycle hook
	Mux       *http.ServeMux
}

// NewHTTPServer builds an HTTP server that will begin serving requests
// when the Fx application starts.
func NewHTTPServer(p ServerParams, mux *http.ServeMux, logging *zap.Logger, service *db.ToDoService) {
	c, _ := configuration.GetConfiguration()
	if c.Configuration.Port == "" {
		log.Fatal("Port value not specified")
	}

	var port = c.Configuration.Port
	logging.Info("Server boot on port", zap.String("port", port))
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
			logging.Info("Starting HTTP server on.")
			ln, err := net.Listen("tcp", server.Addr)
			if err != nil {
				logging.Fatal("HTTP server failed: ", zap.String("error", err.Error()))
				return err
			}
			logging.Info("Starting HTTP server at", zap.String("addr", server.Addr))
			go server.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// This runs when the fx.App stops. On SIGINT.
			logging.Warn("Stopping HTTP server...")
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

		app := fx.New(
			// Add logger first and foremost
			fx.Provide(
				logging.ZapLogger,
			),
			fx.Provide(
				db.DatabaseConfig,
			),
			// 2. Modules (Encapsulate related logic)
			db.DatabaseModule, // Provides *sql.DB and manages its lifecycle
			db.ServiceModule,  // Provides *ToDoService
			// Provide the dependencies
			fx.Provide(
				h.NewHandler,
				h.ProvideMux,
			),

			// invoke server to have it run
			fx.Invoke(NewHTTPServer),
		)

		// starting the application
		app.Run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
