package routes

import (
	"net/http"

	"github.com/alp-tahta/warehouse/internal/handler"
)

func RegisterRoutes(mux *http.ServeMux, h *handler.Handler) {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("POST /order", h.CreateOrder)
	mux.HandleFunc("PUT /barcode/{id}", h.UpdateBarcodeStatus)
	// mux.HandleFunc("GET /product/{id}", corsMiddleware(h.GetProduct))
	// mux.HandleFunc("GET /product", corsMiddleware(h.GetProducts))
	// mux.HandleFunc("DELETE /product/{id}", corsMiddleware(h.DeleteProduct))
}
