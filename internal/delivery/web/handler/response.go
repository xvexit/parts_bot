package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"partsBot/internal/delivery/web/handler/middleware"
)

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("writeJSON error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func getUserIDFromRequest(r *http.Request) (int64, error) {
	if id, ok := middleware.GetUserID(r.Context()); ok {
		return id, nil
	}
	// fallback для ручного тестирования без JWT
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		return 0, errors.New("unauthorized")
	}

	id, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
