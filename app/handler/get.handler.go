package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	// Log event

	// connect with mysql database

	// handle get with pagination
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	if page != "" || limit != "" {
		fmt.Println("GET param page:", page)
		fmt.Println("GET param limit:", limit)
		// ... process it, will be the first (only) if multiple were given
		// note: if they pass in like ?param1=&param2= param1 will also be "" :|
	}

	log.Printf("Received request from %s for path: %s", r.RemoteAddr, r.URL.Path)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "GET Hello from the Cobra-CLI server! Current time: %s", time.Now().Format(time.RFC1123))
}
