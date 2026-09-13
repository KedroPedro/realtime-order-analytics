package order

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	orderuc "github.com/KedroPedro/realtime-order-analytics/internal/application/usecases/order"
	"github.com/rs/zerolog/log"
)

type OrdersHandler struct {
	createUsecase *orderuc.CreateOrderUsecase
}

func NewOrdersHandler(uc *orderuc.CreateOrderUsecase) *OrdersHandler {
	return &OrdersHandler{
		createUsecase: uc,
	}
}

func (h *OrdersHandler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var createOrderReq CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&createOrderReq); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		log.Err(err).Send()
		return
	}

	order, err := createOrderReq.ToEntity()
	if err != nil {
		http.Error(w, "incorrect fields", http.StatusInternalServerError)
		log.Err(err).Send()
		return
	}

	if err := h.createUsecase.Execute(ctx, order); err != nil {
		http.Error(w, "create order error", http.StatusInternalServerError)
		log.Err(err).Send()
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
