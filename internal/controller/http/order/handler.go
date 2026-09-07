package order

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/KedroPedro/realtime-order-analytics/internal/application/usecases"
	"github.com/rs/zerolog/log"
)

type OrdersHandler struct {
	createUsecase *usecases.CreateOrderUsecase
}

func NewOrdersHandler(uc *usecases.CreateOrderUsecase) *OrdersHandler {
	return &OrdersHandler{
		createUsecase: uc,
	}
}

func (h *OrdersHandler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var createOrderReq CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&createOrderReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Err(err).Send()
		return
	}

	order, err := createOrderReq.ToEntity()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Err(err).Send()
		return
	}

	if err := h.createUsecase.Execute(ctx, order); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Err(err).Send()
		return
	}

	w.WriteHeader(http.StatusOK)
}
