package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/leandrobraga/goexpert-desafio-cleanarch/internal/entity"
	"github.com/leandrobraga/goexpert-desafio-cleanarch/internal/pkg/events"
	"github.com/leandrobraga/goexpert-desafio-cleanarch/internal/usecase"
)

type WebOrderHandler struct {
	OrderRepository   entity.OrderRepositoryInterface
	OrderCreatedEvent events.EventInterface
	EventDispatcher   events.EventDispatcherInterface
}

func NewWebOrderHandler(
	OrderRepository entity.OrderRepositoryInterface,
	OrderCreatedEvent events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *WebOrderHandler {
	return &WebOrderHandler{
		OrderRepository:   OrderRepository,
		OrderCreatedEvent: OrderCreatedEvent,
		EventDispatcher:   EventDispatcher,
	}
}

func (h *WebOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createOrder := usecase.NewCreateOrderUseCase(h.OrderRepository, h.OrderCreatedEvent, h.EventDispatcher)
	output, err := createOrder.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *WebOrderHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	dto := usecase.ListOrdersInputDTO{
		Limit: limit,
		Page:  page,
		Sort:  r.URL.Query().Get("page"),
	}
	listOrders := usecase.NewListOrdersUseCase(h.OrderRepository)
	output, err := listOrders.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
