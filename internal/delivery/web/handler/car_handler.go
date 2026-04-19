package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	caruc "partsBot/internal/usecase/car"
	useruc "partsBot/internal/usecase/user"
)

type CarHandler struct {
	carService  *caruc.Service
	userService *useruc.Service
}

func NewCarHandler(carService *caruc.Service, userService *useruc.Service) *CarHandler {
	return &CarHandler{
		carService:  carService,
		userService: userService,
	}
}

// AddCar — POST /api/user/cars
func (h *CarHandler) AddCar(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var dto CarDto

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if strings.TrimSpace(dto.Name) == "" || strings.TrimSpace(dto.VIN) == "" {
		writeError(w, http.StatusBadRequest, "Name and VIN are required")
		return
	}

	input := caruc.NewCarInput(
		dto.Name,
		dto.VIN,
		userID,
	)

	car, err := h.carService.Add(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, car)
}

// MyCars — GET /api/user/cars
func (h *CarHandler) MyCars(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	cars, err := h.carService.ListByUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get cars")
		return
	}

	if len(cars) == 0 {
		writeJSON(w, http.StatusOK, "У вас пока нет добавленных автомобилей")
		return
	}

	writeJSON(w, http.StatusOK, cars)
}
