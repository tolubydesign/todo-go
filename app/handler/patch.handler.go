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

func (h *Handler) PatchHandler(w http.ResponseWriter, r *http.Request) {
	// Log event
	h.logging.Info("Patch Request")

	// Set the Content-Type header for JSON body
	r.Header.Set("Content-Type", "application/json")

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

	if requestBody.Todos != nil {
		for _, todo := range requestBody.Todos {
			opCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			h.logging.Info("updating todo", zap.String("task", todo.Task), zap.String("description", todo.Description))
			// TODO: check if task id is provided

			err := h.service.UpdateToDo(opCtx, todo.ID, &todo.Task, &todo.Description, &todo.Due_date)
			if err != nil {
				h.logging.Warn("update todo failure", zap.Int("id", todo.ID), zap.String("task", todo.Task))
				// TODO: better error handler
				msg := fmt.Sprintf("Failure to update a todo: %d: %s", todo.ID, todo.Task)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
		}
	}

	// Send a response back to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"status": "success"}
	json.NewEncoder(w).Encode(response)
}
