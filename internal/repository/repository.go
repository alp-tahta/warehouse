package repository

import "github.com/alp-tahta/warehouse/internal/model"

type RepositoryI interface {
	CreateOrder(req model.CreateOrderRequest) error
	CheckIfBarcodeMarked(id string) (bool, error)
	MarkBarcode(id string) error
	IncreaseShelfOccupancy(barcodeFields model.BarcodeFields) error
	GetShelvesDetails() ([]model.ShelfInformationWithCustomer, error)
}
