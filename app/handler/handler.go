package handler

import (
	"net/http"
	"time"

	"github.com/tolubydesign/todo-go/app/db"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type RequestBodyToDo struct {
	ID          int        `json:"id,omitempty"`
	Task        string     `json:"task,omitempty"`
	Description string     `json:"description,omitempty"`
	Due_date    string     `json:"due_date,omitempty"`
	Created_at  *time.Time `json:"created_at,omitempty"`
}

type RequestBody struct {
	Todos []RequestBodyToDo `json:"todos"`
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
	mux := http.NewServeMux()

	// Register the handler method
	mux.HandleFunc("GET /todos", NewHandler(service, logging).GetHandler)
	mux.HandleFunc("POST /todos", NewHandler(service, logging).PostHandler)
	mux.HandleFunc("PATCH /todos", NewHandler(service, logging).PatchHandler)
	mux.HandleFunc("DELETE /todos", NewHandler(service, logging).DeleteHandler)
	return mux
}
