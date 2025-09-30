package handler

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func (h *Handler) PatchHandler(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("This is my home page"))
	// Log event

	h.logging.Info("Received request from _ for path: _", zap.String("addr", r.RemoteAddr), zap.String("path", r.URL.Path))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "PATCH Hello from the Cobra-CLI server! Current time: %s", time.Now().Format(time.RFC1123))
}
