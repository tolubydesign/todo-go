package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func (h *Handler) PatchHandler(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("This is my home page"))
	// Log event

	// connect with mysql database

	log.Printf("Received request from %s for path: %s", r.RemoteAddr, r.URL.Path)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "PATCH Hello from the Cobra-CLI server! Current time: %s", time.Now().Format(time.RFC1123))
}
