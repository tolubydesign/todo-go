package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	// Log event
	h.logging.Info("GET Request")

	// NOTE. through the Handler and uber-fx (insane) we already have access to logging and the ToDo Service

	// handle get with pagination
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	if page != "" || limit != "" {
		h.logging.Info("GET param page:", zap.String("page", page))
		h.logging.Info("GET param limit:", zap.String("limit", limit))
		// ... process it, will be the first (only) if multiple were given
		// note: if they pass in like ?param1=&param2= param1 will also be "" :|
	}

	log.Printf("Received request from %s for path: %s", r.RemoteAddr, r.URL.Path)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "GET Hello from the Cobra-CLI server! Current time: %s", time.Now().Format(time.RFC1123))
}
