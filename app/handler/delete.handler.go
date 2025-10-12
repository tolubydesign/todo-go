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

func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	// Log event
	h.logging.Info("DELETE Request")

	// Set the Content-Type header for JSON body
	r.Header.Set("Content-Type", "application/json")

	// Read the request body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		msg := "Request body"
		Response(w, "failed", http.StatusBadRequest, &msg, nil)
		return
	}

	h.logging.Info("DELETE: able to read request body bytes", zap.ByteString("body", bodyBytes))
	// Unmarshal the JSON body into something usable
	var requestBody RequestBody
	err = json.Unmarshal(bodyBytes, &requestBody)
	if err != nil {
		msg := "Unable to read request body"
		Response(w, "failed", http.StatusBadRequest, &msg, nil)
		return
	}

	if len(requestBody.Todos) == 0 {
		h.logging.Info("DELETE: no todos were provided in the request body")
		msg := "Request body invalid"
		Response(w, "failed", http.StatusBadRequest, &msg, nil)
		return
	}

	var dt []int
	if requestBody.Todos != nil {
		for _, todo := range requestBody.Todos {
			opCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			var id *int = todo.ID
			h.logging.Info("DELETE: attempting to delete todo", zap.Any("id", id))

			if id == nil || *id < 0 {
				msg := "invalid id provided"
				Response(w, "failed", http.StatusBadRequest, &msg, nil)
				return
			}

			h.logging.Info("DELETE: deleted", zap.Any("id", id))
			err := h.service.RemoveToDo(opCtx, *todo.ID)
			if err != nil {
				h.logging.Warn("DELETE: delete failed", zap.Int("id", *todo.ID))
				msg := fmt.Sprintf("Failure to create a todo '%d'", todo.ID)
				Response(w, "failed", http.StatusInternalServerError, &msg, nil)
				return
			}

			// create a list of deleted todo ids
			dt = append(dt, *id)
		}
	}

	h.logging.Info("DELETE: deleted items", zap.Any("deleted", dt))
	msg := "request handled successfully"
	Response(w, "successful", http.StatusOK, &msg, dt)
}
