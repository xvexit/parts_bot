package handler

import (
	"encoding/json"
	"net/http"
	"partsBot/internal/usecase/user"
)

type UserHandler struct {
	service *user.Service
}

func NewUserHandler(service *user.Service) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var dto UserDto

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeError(w, http.StatusBadRequest, "Error decoding" + err.Error())
		return
	}

	um := user.NewUserInput(dto.Name, dto.Phone, dto.Password, dto.Email)

	u, err := h.service.Register(r.Context(), um)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, u)
}
