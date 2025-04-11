package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/alp-tahta/warehouse/internal/model"
	"github.com/alp-tahta/warehouse/internal/service"
)

type Handler struct {
	l *slog.Logger
	s service.ServiceI
}

func New(l *slog.Logger, s service.ServiceI) *Handler {
	return &Handler{
		l: l,
		s: s,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("OK"))
	if err != nil {
		h.l.Error("Failed to write response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	productReq := new(model.CreateOrderRequest)

	err := json.NewDecoder(r.Body).Decode(productReq)
	if err != nil {
		h.l.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.s.CreateOrder(*productReq)
	if err != nil {
		h.l.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) UpdateBarcodeStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	err := h.s.UpdateBarcodeStatus(idStr)
	if err != nil {
		h.l.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetShelvesDetails(w http.ResponseWriter, r *http.Request) {
	shelves, err := h.s.GetShelvesDetails()
	if err != nil {
		h.l.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(shelves); err != nil {
		h.l.Error("failed to encode response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetShelvesDetailsHTML(w http.ResponseWriter, r *http.Request) {
	t, m, err := h.s.Index()
	if err != nil {
		h.l.Error("Failed to execute template", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, m)
	if err != nil {
		h.l.Error("Failed to execute template", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
