package handler

import (
	"log"
	"net/http"

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
	// mux.Handle("/", &Handler{})
	mux.HandleFunc("GET /todos", HandlerInstance().GetHandler)
	mux.HandleFunc("POST /todos", HandlerInstance().PostHandler)
	mux.HandleFunc("PATCH /todos", HandlerInstance().PatchHandler)

	log.Println("HTTP handlers registered.")
	return mux
}
