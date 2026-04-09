package httpserver

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// Server is the HTTP transport layer for the mail API (step 1: health only).
type Server struct {
	db *sql.DB
}

// New returns an http.Handler with routes registered.
func New(db *sql.DB) http.Handler {
	s := &Server{db: db}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/api/v1/health", s.handleHealth)
	return mux
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method_not_allowed","message":"GET only"}`, http.StatusMethodNotAllowed)
		return
	}
	if err := s.db.Ping(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"ok":    false,
			"error": "db_unavailable",
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
