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

	var data RequestBody
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		// handle error
		h.logging.Warn("POST: ERROR. shit")

		// Handle error, but if it's an EOF error (no body), it's optional
		if err == io.EOF {
			h.logging.Warn("POST: ERROR. no request body provided", zap.Any("body", r.Body))
			msg := "No request body provided"
			Response(w, "failed", http.StatusBadRequest, &msg, nil)
			return
		} else {
			// Other decoding error, handle as a bad request
			h.logging.Warn("POST: ERROR. invalid request", zap.Any("body", r.Body))
			msg := "Invalid body"
			Response(w, "failed", http.StatusBadRequest, &msg, nil)
			return
		}
	}

	// var body RequestBody
	// if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
	// 	// Handle error, but if it's an EOF error (no body), it's optional
	// 	if err == io.EOF {
	// 		h.logging.Warn("POST: ERROR. no request body provided", zap.Any("body", r.Body))
	// 		msg := "No request body provided"
	// 		Response(w, "failed", http.StatusBadRequest, &msg, nil)
	// 		return
	// 	} else {
	// 		// Other decoding error, handle as a bad request
	// 		h.logging.Warn("POST: ERROR. invalid request", zap.Any("body", r.Body))
	// 		msg := "Invalid body"
	// 		Response(w, "failed", http.StatusBadRequest, &msg, nil)
	// 		return
	// 	}
	// }

	// // Now, check if 'body' is nil
	// if body != nil {
	// 	// Request body was provided and successfully decoded
	// 	fmt.Printf("Received request body: %+v\n", *body)
	// 	// Process the data in 'body'
	// } else {
	// 	// Handle the case where the body was optional and not sent
	// 	fmt.Println("Request body was optional and not sent.")
	// }

	// // Read the request body
	// bodyBytes, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	msg := "Invalid body"
	// 	Response(w, "failed", http.StatusBadRequest, &msg, nil)
	// 	return
	// }
	// defer r.Body.Close()

	// h.logging.Info("POST: able to read request body bytes", zap.ByteString("bytes", bodyBytes))
	// // Unmarshal the JSON body into something usable
	// var requestBody RequestBody
	// err = json.Unmarshal(bodyBytes, &requestBody)
	// if err != nil {

	// 	msg := "invalid request body"
	// 	Response(w, "failed", http.StatusBadRequest, &msg, requestBody)
	// 	return
	// }

	if len(data.Todos) == 0 {
		h.logging.Info("POST: no todos were provided in the request body")
		msg := "Request body not provided"
		Response(w, "failed", http.StatusBadRequest, &msg, nil)
		return
	}

	var todos []ResponseBodyToDo
	if data.Todos != nil {
		for _, todo := range data.Todos {
			opCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			h.logging.Info("POST: adding todo", zap.String("task", todo.Task), zap.String("description", *todo.Description), zap.String("due_date", *todo.Due_date))

			// check that provided task name is of type string & exists
			task := &todo.Task
			if task == nil || len(strings.TrimSpace(todo.Task)) == 0 {
				msg := "invalid task name provided"
				Response(w, "failed", http.StatusBadRequest, &msg, todo)
				return
			}

			// check that provided due_date is valid
			var dueDate *time.Time
			if len(strings.TrimSpace(todo.Task)) < 1 {
				tm, err := helper.IsStringUTC(*todo.Due_date, time.RFC3339)
				// time, err := time.Parse(*todo.Due_date, time.RFC3339)
				if err != nil {
					// something went wrong. Can't use the
					h.logging.Warn("POST: due date provided invalid", zap.String("due_date", *todo.Due_date))
					h.logging.Warn("POST: due date provided invalid. related error", zap.Error(err))
				} else {
					// No issues found
					h.logging.Warn("POST: no errors found", zap.Any("time", tm))
					// dueDate = todo.Due_date
					dueDate = &tm
				}
			}

			t, err := h.service.CreateToDo(opCtx, todo.Task, *todo.Description, dueDate)
			if err != nil {
				h.logging.Warn("create todo failure", zap.String("task", todo.Task))
				msg := fmt.Sprintf("Failure to create a todo '%s'", todo.Task)
				Response(w, "failed", http.StatusInternalServerError, &msg, nil)
				return
			}

			ts := helper.ConvertTimeToString(*t.Due_date)
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
