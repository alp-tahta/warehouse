package service

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"text/template"

	"github.com/alp-tahta/warehouse/internal/barcode"
	"github.com/alp-tahta/warehouse/internal/model"
	"github.com/alp-tahta/warehouse/internal/repository"
)

type Service struct {
	l *slog.Logger
	b barcode.Barcoder
	r repository.RepositoryI
}

func New(l *slog.Logger, b barcode.Barcoder, r repository.RepositoryI) *Service {
	return &Service{
		l: l,
		b: b,
		r: r,
	}
}

type ServiceI interface {
	CreateOrder(req model.CreateOrderRequest) error
	UpdateBarcodeStatus(id string) error
	GetShelvesDetails() ([]model.ShelfInformationWithCustomer, error)
	Index() (*template.Template, []model.ShelfInformationWithCustomer, error)
}

func (s *Service) CreateOrder(req model.CreateOrderRequest) error {
	err := s.r.CreateOrder(req)
	if err != nil {
		s.l.Error("failed to create order", "error", err)
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

func (s *Service) UpdateBarcodeStatus(id string) error {
	marked, err := s.r.CheckIfBarcodeMarked(id)
	if err != nil {
		s.l.Error("failed to check barcode status", "id", id, "error", err)
		return fmt.Errorf("failed to check barcode status: %w", err)
	}

	if marked {
		s.l.Info("barcode already marked", "id", id)
		return fmt.Errorf("barcode already marked: %s", id)
	}

	v, err := s.b.ResolveBarcode(id)
	if err != nil {
		s.l.Error("failed to resolve barcode", "id", id, "error", err)
		return fmt.Errorf("failed to resolve barcode: %w", err)
	}

	err = s.r.IncreaseShelfOccupancy(v)
	if err != nil {
		s.l.Error("failed to increase shelf occupancy", "barcodeFields", v, "error", err)
		return fmt.Errorf("failed to increase shelf occupancy: %w", err)
	}

	err = s.r.MarkBarcode(id)
	if err != nil {
		s.l.Error("failed to mark barcode", "id", id, "error", err)
		return fmt.Errorf("failed to mark barcode: %w", err)
	}

	s.l.Info("barcode marked successfully", "id", id)
	return nil
}

func (s *Service) GetShelvesDetails() ([]model.ShelfInformationWithCustomer, error) {
	s.l.Info("getting shelves details")

	shelves, err := s.r.GetShelvesDetails()
	if err != nil {
		s.l.Error("failed to get shelves details", "error", err)
		return nil, fmt.Errorf("failed to get shelves details: %w", err)
	}

	s.l.Info("shelves details retrieved successfully", "count", len(shelves))
	return shelves, nil
}

func (s *Service) Index() (*template.Template, []model.ShelfInformationWithCustomer, error) {
	shelves, err := s.GetShelvesDetails()

	templatePath := filepath.Join("internal", "templates", "index.html")
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		s.l.Error("Failed to parse template", "error", err)
		return nil, nil, err
	}
	return t, shelves, nil
}
