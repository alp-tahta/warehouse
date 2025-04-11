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
	db *sql.DB
}

func New(l *slog.Logger, db *sql.DB) *Repository {
	return &Repository{
		l:  l,
		db: db,
	}
}

func (r *Repository) CreateOrder(req model.CreateOrderRequest) error {
	queryOrder := `INSERT INTO orders (customer_id) VALUES ($1) RETURNING id`
	var orderID string
	err := r.db.QueryRow(queryOrder, req.CustomerID).Scan(&orderID)
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

		barcode := barcode.CreateBarcodeString(req.CustomerID, orderID, v.ProductID)
		queryBarcode := `INSERT INTO barcodes (order_id, user_id, product_id, code) VALUES ($1, $2, $3, $4)`
		_, err = r.db.Exec(queryBarcode, orderID, req.CustomerID, v.ProductID, barcode)
		if err != nil {
			return fmt.Errorf("could not insert to barcodes: %w", err)
		}
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
	query := `UPDATE barcodes SET marked = true WHERE code = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("could not mark barcode: %w", err)
	}

	return nil
}

// func (r *Repository) GetProduct(id int) (*model.Product, error) {
// 	query := `SELECT id, name, description, price FROM products WHERE id = $1`
// 	product := &model.Product{}
// 	err := r.db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Description, &product.Price)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, fmt.Errorf("product not found")
// 		}
// 		return nil, fmt.Errorf("could not get product: %w", err)
// 	}
// 	return product, nil
// }

// func (r *Repository) GetProducts(ids []int) ([]model.Product, error) {
// 	if len(ids) == 0 {
// 		return []model.Product{}, nil
// 	}

// 	// Create a parameterized query with the correct number of placeholders
// 	query := `SELECT id, name, description, price FROM products WHERE id IN (`
// 	for i := range ids {
// 		if i > 0 {
// 			query += ","
// 		}
// 		query += fmt.Sprintf("$%d", i+1)
// 	}
// 	query += ")"

// 	// Convert ids slice to interface{} slice for Exec
// 	args := make([]interface{}, len(ids))
// 	for i, id := range ids {
// 		args[i] = id
// 	}

// 	// Execute the query
// 	rows, err := r.db.Query(query, args...)
// 	if err != nil {
// 		return nil, fmt.Errorf("could not query products: %w", err)
// 	}
// 	defer rows.Close()

// 	// Process the results
// 	var products []model.Product
// 	for rows.Next() {
// 		var product model.Product
// 		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price)
// 		if err != nil {
// 			return nil, fmt.Errorf("could not scan product: %w", err)
// 		}
// 		products = append(products, product)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("error iterating product rows: %w", err)
// 	}

// 	return products, nil
// }

// func (r *Repository) DeleteProduct(id int) error {
// 	query := `DELETE FROM products WHERE id = $1`
// 	result, err := r.db.Exec(query, id)
// 	if err != nil {
// 		return fmt.Errorf("could not delete product: %w", err)
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return fmt.Errorf("could not get rows affected: %w", err)
// 	}

// 	if rowsAffected == 0 {
// 		return fmt.Errorf("product not found")
// 	}

// 	return nil
// }
