package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	caruc "partsBot/internal/usecase/car"
	useruc "partsBot/internal/usecase/user"

	"github.com/gorilla/mux"
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

	writeJSON(w, http.StatusCreated, serializeCar(car))
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

	out := make([]carResponse, 0, len(cars))
	for i := range cars {
		out = append(out, serializeCar(&cars[i]))
	}
	writeJSON(w, http.StatusOK, out)
}

// DeleteCar — DELETE /api/user/cars/{id}
func (h *CarHandler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	idRaw := mux.Vars(r)["id"]
	carID, err := strconv.ParseInt(strings.TrimSpace(idRaw), 10, 64)
	if err != nil || carID <= 0 {
		writeError(w, http.StatusBadRequest, "Invalid car id")
		return
	}

	car, err := h.carService.GetByID(r.Context(), carID)
	if err != nil {
		writeError(w, http.StatusNotFound, "Автомобиль не найден")
		return
	}

	if car.UserId() != userID {
		writeError(w, http.StatusForbidden, "Нет доступа к удалению этого автомобиля")
		return
	}

	if err := h.carService.Delete(r.Context(), carID); err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось удалить автомобиль")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Автомобиль удален"})
}

// carResponse — доменная car.Car с неэкспортируемыми полями в JSON не попадает; отдаём явный DTO.
type carResponse struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Vin    string `json:"vin"`
	UserID int64  `json:"user_id"`
}

type carSerializer interface {
	ID() int64
	Name() string
	Vin() string
	UserId() int64
}

func serializeCar(c carSerializer) carResponse {
	return carResponse{
		ID:     c.ID(),
		Name:   c.Name(),
		Vin:    c.Vin(),
		UserID: c.UserId(),
	}
}
