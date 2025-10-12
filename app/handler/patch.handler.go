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
	h.logging.Info("PATCH Request")

	// Set the Content-Type header for JSON body
	r.Header.Set("Content-Type", "application/json")

	// Read the request body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.logging.Warn("PATCH: error reading body:", zap.Error(err))
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Unmarshal the JSON body into something usable
	var requestBody RequestBody
	err = json.Unmarshal(bodyBytes, &requestBody)
	if err != nil {
		h.logging.Warn("PATCH: error unmarshal'ing request body:", zap.Error(err))
		http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
		return
	}

	if requestBody.Todos != nil {
		for _, todo := range requestBody.Todos {
			opCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			if todo.ID == nil || *todo.ID == 0 {
				h.logging.Warn("PATCH: Invalid id provided", zap.Int("id", *todo.ID))
				msg := fmt.Sprintf("Invalid ID provided: %d: '%s'", todo.ID, todo.Task)
				Response(w, "failed", http.StatusBadRequest, &msg, nil)
				return
			}

			var id any = todo.ID
			h.logging.Info("PATCH: todo id", zap.Int("id", *todo.ID))
			h.logging.Info("PATCH: id any", zap.Any("id", id))
			// NOTE. confirm that value provided is of type (*int) OR (int)
			pInt, ok := id.(*int)
			idInt, k := id.(int)
			if ok || k {
				h.logging.Info("PATCH: confirmed that id is of type '*int' or 'int'")
				h.logging.Info("PATCH: id provided:", zap.Any("(*int)", pInt), zap.Any("(int)", idInt))

				err := h.service.UpdateToDo(opCtx, todo.ID, &todo.Task, todo.Description, todo.Due_date)
				if err != nil {
					h.logging.Warn("PATCH: update todo failure", zap.Int("id", *todo.ID), zap.String("task", todo.Task))
					msg := fmt.Sprintf("Failed to update a todo: %d: %s", todo.ID, todo.Task)
					Response(w, "failed", http.StatusBadRequest, &msg, nil)
					return
				}
			} else {
				h.logging.Warn("PATCH: Invalid todo id type provided", zap.Int("id", *todo.ID))
				msg := fmt.Sprintf("Invalid id provided: %d: %s", todo.ID, todo.Task)
				Response(w, "failed", http.StatusBadRequest, &msg, nil)
				return
			}

		}
	}

	msg := "request handled successfully"
	Response(w, "successful", http.StatusOK, &msg, nil)
}
