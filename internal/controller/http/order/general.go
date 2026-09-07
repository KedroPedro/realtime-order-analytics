package order

import (
	"net/http"

	"github.com/KedroPedro/realtime-order-analytics/internal/application/usecases"
)

func SetupOrdersRoute(mux *http.ServeMux, uc *usecases.CreateOrderUsecase) {
	handler := NewOrdersHandler(uc)

	mux.HandleFunc("POST /orders/create", handler.HandleCreateOrder)
}
