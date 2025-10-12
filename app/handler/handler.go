package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/tolubydesign/todo-go/app/db"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type RequestBodyToDo struct {
	ID *int `json:"id,omitempty"`

	// max length: 255
	Task        string  `json:"task,omitempty"`
	Description *string `json:"description,omitempty"`

	// RFC 3339. Example 1985-04-12T23:20:50.52Z | 1996-12-19T16:39:57-08:00
	Due_date *string `json:"due_date,omitempty"`

	// Created by database. Example 2025-10-01 21:32:50
	Created_at *time.Time `json:"created_at,omitempty"`
}

type RequestBody struct {
	Todos []RequestBodyToDo `json:"todos"`
}

type ReturnResponse struct {
	// Must be string of "successful" | "failed"
	Status  string `json:"status,omitempty"`
	Message string `json:"message"`
	Data    any    `json:"data"`
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

/*
Setup and handle http response

status: only type `successful` for `failed`
*/
func Response(w http.ResponseWriter, status string, code int, message *string, data any) {
	var res ReturnResponse
	var msg *string

	if message != nil && len(strings.TrimSpace(*message)) > 0 {
		msg = message
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	switch status {
	case "successful":
		res = ReturnResponse{
			Status:  status,
			Message: *msg,
			Data:    data,
		}
	case "failed":
		res = ReturnResponse{
			Status:  status,
			Message: *msg,
			Data:    nil,
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		res = ReturnResponse{
			Status:  status,
			Message: "something has gone wrong",
			Data:    nil,
		}
	}

	json.NewEncoder(w).Encode(res)
}
