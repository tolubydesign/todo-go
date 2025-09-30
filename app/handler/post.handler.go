package handler

import (
	"context"
	"encoding/json"
	"fmt"
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
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received message: %v\n", requestBody.Todos)
	if requestBody.Todos != nil {
		for _, todo := range requestBody.Todos {
			opCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			h.logging.Info("adding todo", zap.String("task", todo.Task), zap.String("description", todo.Description))
			_, err := h.service.CreateToDo(opCtx, todo.Task, todo.Description)
			if err != nil {
				h.logging.Warn("create todo failure", zap.String("task", todo.Task))
				// Todo better error handler
				http.Error(w, "Failure to create a todo", http.StatusBadRequest)
				return
			}
		}
	}

	// Send a response back to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"status": "success", "received_message": fmt.Sprintf("%d", len(requestBody.Todos))}
	json.NewEncoder(w).Encode(response)
}
