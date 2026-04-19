package handler

import (
	"net/http"
	"strings"

	adap "partsBot/internal/infrastructure/adapter"
	caruc "partsBot/internal/usecase/car"
	useruc "partsBot/internal/usecase/user"
)

type PartsHandler struct {
	partsAdapter *adap.Adapter
	carService   *caruc.Service
	userService  *useruc.Service
}

func NewPartsHandler(parts *adap.Adapter, carService *caruc.Service, userService *useruc.Service) *PartsHandler {
	return &PartsHandler{
		partsAdapter: parts,
		carService:   carService,
		userService:  userService,
	}
}

// SearchParts — GET /api/parts/search?part=...
func (h *PartsHandler) SearchParts(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("part"))
	if query == "" {
		writeJSON(w, http.StatusOK, "Query is required")
		return
	}

	userID, err := getUserIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	cars, err := h.carService.ListByUser(r.Context(), userID)
	if err != nil || len(cars) == 0 {
		writeError(w, http.StatusBadRequest, "У пользователя нет добавленных автомобилей")
		return
	}

	vin := cars[0].Vin()

	parts, err := h.partsAdapter.GetPartsByVIN(r.Context(), vin, query, false)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка поиска запчастей")
		return
	}

	if len(parts) == 0 {
		writeJSON(w, http.StatusOK, "По запросу ничего не найдено")
		return
	}

	writeJSON(w, http.StatusOK, parts)
}
