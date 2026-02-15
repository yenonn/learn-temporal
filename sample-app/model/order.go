package model

type Item struct {
	SKU      string  `json:"sku"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

type Order struct {
	OrderID    string  `json:"order_id"`
	CustomerID string  `json:"customer_id"`
	Email      string  `json:"email"`
	Items      []Item  `json:"items"`
	Address    Address `json:"address"`
	TotalAmount float64 `json:"total_amount"`
}

type PaymentResult struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}

type ShipmentResult struct {
	TrackingNumber string `json:"tracking_number"`
	Carrier        string `json:"carrier"`
}

type OrderResult struct {
	OrderID        string         `json:"order_id"`
	Status         string         `json:"status"`
	Payment        PaymentResult  `json:"payment"`
	Shipment       ShipmentResult `json:"shipment"`
}
