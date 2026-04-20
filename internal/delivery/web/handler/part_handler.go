package handler

import (
	"log"
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
	partID := strings.TrimSpace(r.URL.Query().Get("part"))
	if partID == "" {
		writeError(w, http.StatusBadRequest, "part is required")
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

	parts, err := h.partsAdapter.GetPartsByVIN(r.Context(), vin, partID, false)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка поиска запчастей")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items":   parts,
		"message": "ok",
	})
}

// SearchTree — GET /api/user/parts/tree
func (h *PartsHandler) SearchTree(w http.ResponseWriter, r *http.Request) {
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
	nodes, err := h.partsAdapter.GetSearchTreeByVIN(r.Context(), vin)
	if err != nil {
		log.Printf("SearchTree failed for user=%d vin=%s: %v", userID, vin, err)
		writeError(w, http.StatusBadRequest, "Ошибка загрузки дерева запчастей: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items":   nodes,
		"message": "ok",
	})
}
