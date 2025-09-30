package handler

import (
	"log"
	"net/http"

	"github.com/tolubydesign/todo-go/app/db"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type RequestBody struct {
	Todos []db.ToDo `json:"todo"`
}

type Handler struct {
	service *db.ToDoService
	logging *zap.Logger
}

// Handler instance.
func NewHandler(service *db.ToDoService, logging *zap.Logger) *Handler {
	return &Handler{
		service: service,
		logging: logging,
	}
}

// MuxParams defines the dependencies for the HTTP router.
type MuxParams struct {
	fx.In   // Embed fx.In to receive dependencies
	Handler *Handler
}

// ProvideMux registers the HTTP handler and returns a *http.ServeMux.
func ProvideMux(p MuxParams, service *db.ToDoService, logging *zap.Logger) *http.ServeMux {
	log.Println("func ProvideMux.")
	mux := http.NewServeMux()

	// Register the handler method
	// mux.Handle("/", &Handler{})
	mux.HandleFunc("GET /todos", NewHandler(service, logging).GetHandler)
	mux.HandleFunc("POST /todos", NewHandler(service, logging).PostHandler)
	mux.HandleFunc("PATCH /todos", NewHandler(service, logging).PatchHandler)

	log.Println("HTTP handlers registered.")
	return mux
}
