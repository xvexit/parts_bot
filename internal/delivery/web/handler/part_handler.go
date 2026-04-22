package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	adap "partsBot/internal/infrastructure/adapter"
	caruc "partsBot/internal/usecase/car"
	partuc "partsBot/internal/usecase/part"
	useruc "partsBot/internal/usecase/user"
)

type PartsHandler struct {
	partsAdapter *adap.Adapter
	carService   *caruc.Service
	partService  *partuc.Service
	userService  *useruc.Service
}

func NewPartsHandler(parts *adap.Adapter, carService *caruc.Service, partService *partuc.Service, userService *useruc.Service) *PartsHandler {
	return &PartsHandler{
		partsAdapter: parts,
		carService:   carService,
		partService:  partService,
		userService:  userService,
	}
}

// SearchParts — GET /api/parts/search?q=...
func (h *PartsHandler) SearchParts(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		writeError(w, http.StatusBadRequest, "q is required")
		return
	}

	items, err := h.partService.Search(r.Context(), query, 30)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка поиска по каталогу")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items":   toPartResponse(items),
		"message": "Каталог обновлен",
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

// CheckPartOffer — POST /api/user/parts/check
// Нейтральная "уточнялка": если внешний сервис не ответил, UI не ломаем.
func (h *PartsHandler) CheckPartOffer(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req struct {
		PartID string `json:"part_id"`
		Name   string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	cars, err := h.carService.ListByUser(r.Context(), userID)
	if err != nil || len(cars) == 0 {
		writeJSON(w, http.StatusOK, map[string]any{
			"found":   false,
			"message": "Добавьте автомобиль для уточнения сроков",
		})
		return
	}

	vin := cars[0].Vin()
	parts, err := h.partsAdapter.GetPartsByVIN(r.Context(), vin, "", false)
	if err != nil {
		log.Printf("CheckPartOffer failed user=%d vin=%s: %v", userID, vin, err)
		writeJSON(w, http.StatusOK, map[string]any{
			"found":   false,
			"message": "Сейчас нельзя уточнить сроки, попробуйте позже",
		})
		return
	}

	partIDNeedle := strings.ToLower(strings.TrimSpace(req.PartID))
	nameNeedle := strings.ToLower(strings.TrimSpace(req.Name))
	for _, p := range parts {
		blob := strings.ToLower(p.Name + " " + p.Shortname + " " + p.Parts)
		if (partIDNeedle != "" && strings.Contains(blob, partIDNeedle)) ||
			(nameNeedle != "" && strings.Contains(blob, nameNeedle)) {
			writeJSON(w, http.StatusOK, map[string]any{
				"found":   true,
				"message": "Есть актуальные предложения, срок подтвержден",
			})
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"found":   false,
		"message": "Точное предложение не найдено, оставьте заявку менеджеру",
	})
}

type partResponse struct {
	PartID      string `json:"part_id"`
	Name        string `json:"name"`
	Brand       string `json:"brand"`
	Price       int64  `json:"price"`
	DeliveryDay int    `json:"delivery_day"`
}

func toPartResponse(items []partuc.CatalogItem) []partResponse {
	out := make([]partResponse, 0, len(items))
	for _, item := range items {
		out = append(out, partResponse{
			PartID:      item.PartID,
			Name:        item.Name,
			Brand:       item.Brand,
			Price:       item.Price,
			DeliveryDay: item.DeliveryDay,
		})
	}
	return out
}
