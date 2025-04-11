package handler

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alp-tahta/warehouse/internal/model"
	"github.com/alp-tahta/warehouse/internal/service"
	"github.com/golang/mock/gomock"
)

func TestHealth(t *testing.T) {
	h := New(nil, nil) // Initialize handler with nil logger and service for this test

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	h.Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != "OK" {
		t.Errorf("expected body %q, got %q", "OK", w.Body.String())
	}
}

func TestCreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockServiceI(ctrl)
	logger := slog.New(slog.NewTextHandler(nil, nil))
	h := New(logger, mockService)

	t.Run("valid request", func(t *testing.T) {
		mockService.EXPECT().CreateOrder(gomock.AssignableToTypeOf(model.CreateOrderRequest{})).Return(nil)

		order := model.CreateOrderRequest{
			CustomerID: "123",
			OrderItems: []model.CreateOrderItemRequest{
				{
					ProductID: 456,
					Quantity:  2,
				},
			},
		}
		body, _ := json.Marshal(order)
		req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.CreateOrder(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}
	})
}
