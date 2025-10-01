package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	// NOTE. through the Handler and uber-fx (insane) we already have access to logging and the ToDo Service
	h.logging.Info("GET Request")

	// Set the Content-Type header for JSON body
	r.Header.Set("Content-Type", "application/json")

	// set default variables
	limit := "10"
	// NOTE. page number is incremented by 10. so page 1 would be 10 items in. page 2 would be 20 items in
	page := "1"

	// set the variables
	page = r.URL.Query().Get("page")
	limit = r.URL.Query().Get("limit")

	// validate that page and limit are both int dressed as string
	limitNum, err := strconv.Atoi(limit) // Atoi returns an int and an error
	if err != nil {
		h.logging.Info("limit param is not a number:", zap.String("limit", limit))
		// default to
		limit = "10"
		limitNum = 10
	}

	pageNum, err := strconv.Atoi(page) // Atoi returns an int and an error
	if err != nil {
		h.logging.Info("page params is not a number:", zap.String("page", page))
		// default to
		page = "0" // = 0
		pageNum = 0
	}

	// TODO: (future development) increment based on "limit"
	if pageNum > 0 {
		pageNum = (pageNum - 1) * 10 // increment by 10 per page viewed
	}

	// The above section is scuffed I know

	// convert back to string
	pageStr := strconv.Itoa(pageNum)
	limitStr := strconv.Itoa(limitNum)

	h.logging.Info("GET param page:", zap.String("page", page), zap.Int("num", pageNum), zap.String("str", pageStr))
	h.logging.Info("GET param limit:", zap.String("limit", limit), zap.Int("num", limitNum), zap.String("str", limitStr))

	var responseToDos []RequestBodyToDo
	opCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// handle get with pagination in query
	todos, err := h.service.GetToDo(opCtx, limitStr, pageStr)
	if err != nil {
		h.logging.Warn("get todo failure")
		// TODO: better error handler
		http.Error(w, "Failure to get a todos", http.StatusInternalServerError)
		return
	}

	for _, todo := range todos {
		responseToDos = append(responseToDos, RequestBodyToDo{
			ID:          todo.ID,
			Task:        todo.Task,
			Description: todo.Task_description,
			// Due_date: time.Now().Format(todo.Due_date), // needs work
		})
	}

	// // Marshal the slice into a JSON byte slice
	// marshalTodos, err := json.Marshal(responseToDos)
	// if err != nil {
	// 	http.Error(w, "Failure to return added todos", http.StatusInternalServerError)
	// 	return
	// }
	// h.logging.Info("List of todos: ", zap.ByteString("todos", marshalTodos))

	// Send a response back to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string][]RequestBodyToDo{"todos": responseToDos}
	json.NewEncoder(w).Encode(response)
}
