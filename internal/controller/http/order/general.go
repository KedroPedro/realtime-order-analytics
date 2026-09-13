package order

import (
	"net/http"

	orderuc "github.com/KedroPedro/realtime-order-analytics/internal/application/usecases/order"
)

func SetupOrdersRoute(mux *http.ServeMux, uc *orderuc.CreateOrderUsecase) {
	handler := NewOrdersHandler(uc)

	mux.HandleFunc("POST /orders/create", handler.HandleCreateOrder)
}
