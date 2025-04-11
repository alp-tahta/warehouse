package routes

import (
	"net/http"

	"github.com/alp-tahta/warehouse/internal/handler"
)

func RegisterRoutes(mux *http.ServeMux, h *handler.Handler) {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("POST /order", h.CreateOrder)
	mux.HandleFunc("PUT /barcode/{id}", h.UpdateBarcodeStatus)
	mux.HandleFunc("GET /shelf", h.GetShelvesDetails)
	mux.HandleFunc("GET /shelf-html", h.GetShelvesDetailsHTML)
}
