/*
Copyright © 2025 Tolu Adesina <tolu.adesina...>
*/
package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"net/http"

	"github.com/spf13/cobra"
	configuration "github.com/tolubydesign/todo-go/app/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type EchoHandler struct {
	log *zap.Logger
}

func NewEchoHandler(log *zap.Logger) *EchoHandler {
	return &EchoHandler{log: log}
}

type Handler struct{}

// Handler instance.
func HandlerInstance() *Handler {
	return &Handler{}
}

// MuxParams defines the dependencies for the HTTP router.
type MuxParams struct {
	fx.In   // Embed fx.In to receive dependencies
	Handler *Handler
}

// ProvideMux registers the HTTP handler and returns a *http.ServeMux.
func ProvideMux(p MuxParams) *http.ServeMux {
	log.Println("func ProvideMux.")
	mux := http.NewServeMux()

	// Register the handler method
	mux.Handle("/", &Handler{})
	mux.HandleFunc("GET /todos", p.Handler.GetToDoHandler)
	mux.HandleFunc("POST /todos", p.Handler.PostToDoHandler)
	mux.HandleFunc("PATCH /todos", p.Handler.PatchToDoHandler)

	log.Println("HTTP handlers registered.")
	return mux
}

func (h *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := io.Copy(w, r.Body); err != nil {
		h.log.Warn("Failed to handle request", zap.Error(err))
	}
}

// ServerParams defines the dependencies for the HTTP server setup.
type ServerParams struct {
	fx.In
	Lifecycle fx.Lifecycle   // The fx lifecycle hook
	Mux       *http.ServeMux // The configured router
}

func (h *Handler) GetToDoHandler(w http.ResponseWriter, r *http.Request) {
	// Log event

	// connect with mysql database

	// handle get with pagination
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	if page != "" || limit != "" {
		fmt.Println("GET param page:", page)
		fmt.Println("GET param limit:", limit)
		// ... process it, will be the first (only) if multiple were given
		// note: if they pass in like ?param1=&param2= param1 will also be "" :|
	}

	log.Printf("Received request from %s for path: %s", r.RemoteAddr, r.URL.Path)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "GET Hello from the Cobra-CLI server! Current time: %s", time.Now().Format(time.RFC1123))
}

func (h *Handler) PostToDoHandler(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("This is my home page"))
	// Log event

	// connect with mysql database

	log.Printf("Received request from %s for path: %s", r.RemoteAddr, r.URL.Path)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "POST Hello from the Cobra-CLI server! Current time: %s", time.Now().Format(time.RFC1123))
}

func (h *Handler) PatchToDoHandler(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("This is my home page"))
	// Log event

	// connect with mysql database

	log.Printf("Received request from %s for path: %s", r.RemoteAddr, r.URL.Path)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "PATCH Hello from the Cobra-CLI server! Current time: %s", time.Now().Format(time.RFC1123))
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
			// // Use a goroutine to start serving so it doesn't block the OnStart hook.
			// go func() {
			// 	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// 		log.Fatalf("HTTP server failed: %v", err)
			// 	}
			// }()
			// return nil
		},
		OnStop: func(ctx context.Context) error {
			// This runs when the fx.App stops. On SIGINT.
			log.Warn("Stopping HTTP server...")
			return server.Shutdown(ctx)
		},
	})
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the ToDo Backend API",
	Long:  "This ToDo application uses Cobra for CLI commands and uber-fx for dependency injection and lifecycle management.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")

		config, err := configuration.BuildConfiguration()
		if err != nil {
			fmt.Println("error", err.Error())
			panic(err)
		}

		c := config.Configuration
		fmt.Println("ENVIRONMENT :::", c.Environment)
		fmt.Println("Server is running on port ", c.Port)

		app := fx.New(
			// Provide the dependencies
			fx.Provide(
				HandlerInstance,
				ProvideMux,
				NewEchoHandler,
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
