package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func (h *Handler) PostHandler(w http.ResponseWriter, r *http.Request) {
	// Log event
	h.logging.Info("POST Request")

	// Set the Content-Type header for JSON body
	r.Header.Set("Content-Type", "application/json")

	// Already connected to database
	h.logging.Info("Received request from %s for path: %s", zap.String("remote-addr", r.RemoteAddr), zap.String("path", r.URL.Path))

	// Read the request body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	h.logging.Info("able to read body bytes")
	// Unmarshal the JSON body into something usable
	var requestBody RequestBody
	err = json.Unmarshal(bodyBytes, &requestBody)
	if err != nil {
		http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
		return
	}

	var todos []RequestBodyToDo
	if requestBody.Todos != nil {
		for _, todo := range requestBody.Todos {
			opCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			h.logging.Info("adding todo", zap.String("task", todo.Task), zap.String("description", todo.Description))
			todo, err := h.service.CreateToDo(opCtx, todo.Task, todo.Description)
			if err != nil {
				h.logging.Warn("create todo failure", zap.String("task", todo.Task))
				// Todo better error handler
				http.Error(w, "Failure to create a todo", http.StatusBadRequest)
				return
			}

			// add returning todos with ids
			todos = append(todos, RequestBodyToDo{
				ID:          todo.ID,
				Task:        todo.Task,
				Description: todo.Task_description,
			})
		}
	}

	// Marshal the slice into a JSON byte slice
	toDoJSONData, err := json.Marshal(todos)
	if err != nil {
		http.Error(w, "Failure to return added todos", http.StatusInternalServerError)
		return
	}

	h.logging.Info("todo json data", zap.ByteString("todo json data", toDoJSONData))

	// Send a response back to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// NOTE. alternative method of returning data in JSON
	// response := map[string][]RequestBodyToDo{"todos": todos}
	// json.NewEncoder(w).Encode(response)
	//
	w.Write(toDoJSONData)
}
