package service

import (
	"fmt"
	"log/slog"

	"github.com/alp-tahta/warehouse/internal/barcode"
	"github.com/alp-tahta/warehouse/internal/model"
	"github.com/alp-tahta/warehouse/internal/repository"
)

type Service struct {
	l *slog.Logger
	r repository.RepositoryI
}

func New(l *slog.Logger, r repository.RepositoryI) *Service {
	return &Service{
		l: l,
		r: r,
	}
}

type ServiceI interface {
	CreateOrder(req model.CreateOrderRequest) error
	UpdateBarcodeStatus(id string) error
	GetShelvesDetails() ([]model.ShelfInformationWithCustomer, error)
	// GetProduct(id int) (*model.Product, error)
	// GetProducts(ids []int) ([]model.Product, error)
	// DeleteProduct(id int) error
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

	v, err := barcode.ResolveBarcode(id)
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

// func (s *Service) GetProduct(id int) (*model.Product, error) {
// 	s.l.Info("getting product", "id", id)

// 	product, err := s.r.GetProduct(id)
// 	if err != nil {
// 		s.l.Error("failed to get product", "id", id, "error", err)
// 		return nil, fmt.Errorf("failed to get product: %w", err)
// 	}

// 	s.l.Info("product retrieved successfully", "id", id)
// 	return product, nil
// }

// func (s *Service) DeleteProduct(id int) error {
// 	s.l.Info("deleting product", "id", id)

// 	err := s.r.DeleteProduct(id)
// 	if err != nil {
// 		s.l.Error("failed to delete product", "id", id, "error", err)
// 		return fmt.Errorf("failed to delete product: %w", err)
// 	}

// 	s.l.Info("product deleted successfully", "id", id)
// 	return nil
// }

// func (s *Service) GetProducts(ids []int) ([]model.Product, error) {
// 	s.l.Info("getting products", "ids", ids)

// 	products, err := s.r.GetProducts(ids)
// 	if err != nil {
// 		s.l.Error("failed to get products", "ids", ids, "error", err)
// 		return nil, fmt.Errorf("failed to get products: %w", err)
// 	}

// 	s.l.Info("products retrieved successfully", "count", len(products))
// 	return products, nil
// }
