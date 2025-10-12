package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

func (h *Handler) PostHandler(w http.ResponseWriter, r *http.Request) {
	// Log event
	h.logging.Info("POST Request")

	// Set the Content-Type header for JSON body
	r.Header.Set("Content-Type", "application/json")

	// Read the request body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		msg := "Request body"
		Response(w, "failed", http.StatusBadRequest, &msg, nil)
		return
	}

	h.logging.Info("able to read request body bytes")
	// Unmarshal the JSON body into something usable
	var requestBody RequestBody
	err = json.Unmarshal(bodyBytes, &requestBody)
	if err != nil {
		msg := "Request body unmarshaling"
		Response(w, "failed", http.StatusBadRequest, &msg, nil)
		return
	}

	if len(requestBody.Todos) == 0 {
		h.logging.Info("no todos were provided in the request body")
		msg := "Request body not provided"
		Response(w, "failed", http.StatusBadRequest, &msg, nil)
		return
	}

	var todos []RequestBodyToDo
	if requestBody.Todos != nil {
		for _, todo := range requestBody.Todos {
			opCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			h.logging.Info("adding todo", zap.String("task", todo.Task), zap.String("description", *todo.Description))

			task := &todo.Task
			if task == nil || len(strings.TrimSpace(todo.Task)) > 0 {
				msg := "no task name provided"
				Response(w, "failed", http.StatusBadRequest, &msg, nil)
				return
			}

			t, err := h.service.CreateToDo(opCtx, todo.Task, *todo.Description)
			if err != nil {
				h.logging.Warn("create todo failure", zap.String("task", todo.Task))
				msg := fmt.Sprintf("Failure to create a todo '%s'", todo.Task)
				Response(w, "failed", http.StatusInternalServerError, &msg, nil)
				return
			}

			// add returning todos with ids
			todos = append(todos, RequestBodyToDo{
				ID:          &t.ID,
				Task:        t.Task,
				Description: &t.Task_description,
			})
		}
	}

	// Marshal the slice into a JSON byte slice
	todoJSON, err := json.Marshal(todos)
	if err != nil {
		msg := "unable to return added to-dos"
		Response(w, "failed", http.StatusInternalServerError, &msg, nil)
		return
	}

	h.logging.Info("todo json data", zap.ByteString("todo json data", todoJSON))

	msg := "request handled successfully"
	Response(w, "successful", http.StatusOK, &msg, todos)
}
