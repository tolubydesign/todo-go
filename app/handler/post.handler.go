package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/tolubydesign/todo-go/app/helper"
	"go.uber.org/zap"
)

func (h *Handler) PostHandler(w http.ResponseWriter, r *http.Request) {
	// Log event
	h.logging.Info("POST: Request")

	// Set the Content-Type header for JSON body
	r.Header.Set("Content-Type", "application/json")

	// Read the request body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		msg := "Invalid body"
		Response(w, "failed", http.StatusBadRequest, &msg, nil)
		return
	}
	defer r.Body.Close()

	h.logging.Info("POST: able to read request body bytes", zap.ByteString("bytes", bodyBytes))
	// Unmarshal the JSON body into something usable
	var requestBody RequestBody

	err = json.Unmarshal(bodyBytes, &requestBody)
	if err != nil {
		msg := "invalid request body"
		Response(w, "failed", http.StatusBadRequest, &msg, requestBody)
		return
	}

	if len(requestBody.Todos) == 0 {
		h.logging.Warn("POST: no todos were provided in the request body")
		msg := "Request body not provided"
		Response(w, "failed", http.StatusBadRequest, &msg, nil)
		return
	}

	var todos []ResponseBodyToDo
	if requestBody.Todos != nil {
		for i, todo := range requestBody.Todos {
			opCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			h.logging.Info("POST: adding todo", zap.String("task", todo.Task), zap.String("description", todo.Description), zap.String("due_date", todo.Due_date))

			// check that provided task name is of type string & exists
			if len(strings.TrimSpace(todo.Task)) == 0 {
				msg := "invalid task name provided"
				Response(w, "failed", http.StatusBadRequest, &msg, todo)
				return
			}

			// check that provided due_date is valid
			var dueDate *time.Time
			tm, err := helper.IsStringUTC(todo.Due_date, time.RFC3339)
			if err != nil {
				// Something went wrong. Can't use the
				h.logging.Warn("POST: due date provided invalid", zap.String("due_date", todo.Due_date))
				h.logging.Warn("POST: due date provided invalid. related error", zap.Error(err))
			} else {
				// No issues found
				h.logging.Warn("POST: no errors found", zap.Any("time", tm))
				// dueDate = todo.Due_date
				dueDate = &tm
			}

			h.logging.Info("POST: creating todo. task placement:", zap.Int("int", i), zap.String("task", todo.Task), zap.String("description", todo.Description), zap.String("due_date", todo.Due_date))

			t, err := h.service.CreateToDo(opCtx, todo.Task, todo.Description, dueDate)
			if err != nil {
				h.logging.Warn("create todo failure", zap.String("task", todo.Task))
				msg := fmt.Sprintf("Failure to create a todo '%s'", todo.Task)
				Response(w, "failed", http.StatusInternalServerError, &msg, nil)
				return
			}

			var ts string
			if dueDate != nil {
				ts = helper.ConvertTimeToString(*t.Due_date)
			}
			create_ts := helper.ConvertTimeToString(t.Created_at)

			// add returning todos with ids
			todos = append(todos, ResponseBodyToDo{
				ID:          &t.ID,
				Task:        t.Task,
				Description: &t.Task_description,
				Due_date:    &ts,
				Created_at:  create_ts,
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

	msg := helper.SuccessfulResponseMessage
	Response(w, "successful", http.StatusOK, &msg, todos)
}
