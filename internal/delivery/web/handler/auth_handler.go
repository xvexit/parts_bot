package handler

import (
	"encoding/json"
	"net/http"
	"partsBot/internal/usecase/auth"
)

type AuthHandler struct {
	authService *auth.Service
}

func NewAuthHandler(authService *auth.Service) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// LOGIN
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var dto LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	access, refresh, err := h.authService.Login(r.Context(), dto.Email, dto.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

// REFRESH
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var dto RefreshRequest

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	access, err := h.authService.Refresh(r.Context(), dto.RefreshToken)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"access_token": access,
	})
}
