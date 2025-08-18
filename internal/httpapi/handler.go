package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strings"
	"wb-tech-L0/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LISTEN HERE:", r.URL.String())
	route := strings.Split(r.URL.String(), "/")
	// тут внимательно, видимо, первый элемент слайса - пустая строка
	orderUID, err := uuid.Parse(route[len(route)-1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")

	order, err := h.service.GetOrder(r.Context(), orderUID)
	fmt.Printf("ORDER: %v  UID: %v\n", order, orderUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, fmt.Sprintf("invalid id: %v", err), http.StatusNotFound)
		}
	}

	err = json.NewEncoder(w).Encode(order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
