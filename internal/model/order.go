package model

type Order struct {
	ID         int `json:"id"`
	CustomerID int `json:"CustomerID"`
	OrderItems []OrderItem
}

type OrderItem struct {
	ID        int `json:"id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
	OrderID   int `json:"order_id"`
}

type Product struct {
	ID   int
	Name int
}

type Customer struct {
	ID   int
	name string
}

type Shelf struct {
	Location string
}

type CreateOrderRequest struct {
	CustomerID string                   `json:"customer_id"`
	OrderItems []CreateOrderItemRequest `json:"order_items"`
}

type CreateOrderItemRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type BarcodeFields struct {
	CustomerID string
	OrderID    string
	ProductID  int
}

type ShelfInformationWithCustomer struct {
	CustomerID       string
	OrderID          string
	ShelfName        string
	CurrentOccupancy int
	Capacity         int
}
