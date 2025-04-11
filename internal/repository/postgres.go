package repository

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/alp-tahta/warehouse/internal/barcode"
	"github.com/alp-tahta/warehouse/internal/model"
)

type Repository struct {
	l  *slog.Logger
	b  barcode.Barcoder
	db *sql.DB
}

func New(l *slog.Logger, b barcode.Barcoder, db *sql.DB) *Repository {
	return &Repository{
		l:  l,
		b:  b,
		db: db,
	}
}

func (r *Repository) CreateOrder(req model.CreateOrderRequest) error {
	err := r.FreeFullShelves()
	if err != nil {
		return fmt.Errorf("could not free full shelves: %w", err)
	}

	queryOrder := `INSERT INTO orders (customer_id) VALUES ($1) RETURNING id`
	var orderID string
	err = r.db.QueryRow(queryOrder, req.CustomerID).Scan(&orderID)
	if err != nil {
		return fmt.Errorf("could not insert to orders: %w", err)
	}

	// Insert the order items
	for _, v := range req.OrderItems {
		queryItems := `INSERT INTO order_items (product_id, quantity, order_id) VALUES ($1, $2, $3)`
		_, err = r.db.Exec(queryItems, v.ProductID, v.Quantity, orderID)
		if err != nil {
			return fmt.Errorf("could not insert to order items: %w", err)
		}

		barcode := r.b.CreateBarcodeString(req.CustomerID, orderID, v.ProductID)
		queryBarcode := `INSERT INTO barcodes (order_id, user_id, product_id, code) VALUES ($1, $2, $3, $4)`
		_, err = r.db.Exec(queryBarcode, orderID, req.CustomerID, v.ProductID, barcode)
		if err != nil {
			return fmt.Errorf("could not insert to barcodes: %w", err)
		}
	}

	shelf, err := r.FindFreeShelf()
	if err != nil {
		return fmt.Errorf("could not find free shelf: %w", err)
	}

	err = r.SpareShelfForOrder(shelf, req.CustomerID, orderID, len(req.OrderItems))
	if err != nil {
		return fmt.Errorf("could not spare shelf for order: %w", err)
	}

	return nil
}

func (r *Repository) CheckIfBarcodeMarked(id string) (bool, error) {
	query := `SELECT COUNT(*) FROM barcodes WHERE code = $1 AND marked = true`
	var count int
	err := r.db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("could not check barcode: %w", err)
	}

	return count > 0, nil
}

func (r *Repository) MarkBarcode(id string) error {
	// Check if the barcode marked
	queryMarked := `SELECT marked FROM barcodes WHERE code = $1`
	var count bool
	err := r.db.QueryRow(queryMarked, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("could not check barcode: %w", err)
	}
	if count == true {
		return fmt.Errorf("barcode already marked: %s", id)
	}

	query := `UPDATE barcodes SET marked = true WHERE code = $1`
	_, err = r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("could not mark barcode: %w", err)
	}

	return nil
}

// Free full shelves
func (r *Repository) FreeFullShelves() error {
	query := `UPDATE shelves SET user_id = NULL, order_id = NULL, current_occupancy = 0, capacity = 0  WHERE current_occupancy = capacity`
	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("could not free full shelves: %w", err)
	}

	return nil
}

// Find empty shelf
func (r *Repository) FindFreeShelf() (string, error) {
	query := `SELECT name FROM shelves WHERE user_id IS NULL AND order_id IS NULL AND capacity = 0 ORDER BY id LIMIT 1`
	var shelfName string
	err := r.db.QueryRow(query).Scan(&shelfName)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no empty shelf found")
		}
		return "", fmt.Errorf("could not find empty shelf: %w", err)
	}

	return shelfName, nil
}

func (r *Repository) SpareShelfForOrder(shelfName, userID, orderID string, capacity int) error {
	query := `UPDATE shelves SET user_id = $1, order_id = $2, current_occupancy = 0, capacity = $3 WHERE name = $4`
	_, err := r.db.Exec(query, userID, orderID, capacity, shelfName)
	if err != nil {
		return fmt.Errorf("could not spare shelf for order: %w", err)
	}
	return nil
}

// Increase shelf occupancy by 1
func (r *Repository) IncreaseShelfOccupancy(barcodeFields model.BarcodeFields) error {
	query := `UPDATE shelves SET current_occupancy = current_occupancy + 1 WHERE user_id = $1 AND order_id= $2`
	_, err := r.db.Exec(query, barcodeFields.CustomerID, barcodeFields.OrderID)
	if err != nil {
		return fmt.Errorf("could not increase shelf occupancy: %w", err)
	}
	return nil
}

func (r *Repository) GetShelvesDetails() ([]model.ShelfInformationWithCustomer, error) {
	query := `SELECT name, user_id, order_id, current_occupancy, capacity FROM shelves WHERE user_id IS NOT NULL AND order_id IS NOT NULL ORDER BY id`
	rows, err := r.db.Query(query)
	if err != nil {
		r.l.Error("could not query shelves details", "error", err)
		return nil, err
	}
	defer rows.Close()
	var shelves []model.ShelfInformationWithCustomer
	for rows.Next() {
		var shelf model.ShelfInformationWithCustomer
		err := rows.Scan(&shelf.ShelfName, &shelf.CustomerID, &shelf.OrderID, &shelf.CurrentOccupancy, &shelf.Capacity)
		if err != nil {
			r.l.Error("could not scan shelf details", "error", err)
			return nil, err
		}
		shelves = append(shelves, shelf)
	}

	return shelves, nil
}
