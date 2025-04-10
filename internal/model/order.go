package model

type Order struct {
	ID         int
	CustomerID int
	OrderItems []OrderItem
}

type OrderItem struct {
	ID        int
	ProductID int
	Quantity  int
	OrderID   int
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
