package order

import "time"

type Order struct {
	ID          string           `json:"id"`
	CreatedAt   time.Time        `json:"created_at"`
	AccountId   string           `json:"account_id"`
	TotalAmount float64          `json:"total_amount"`
	Products    []OrderedProduct `json:"products"`
}

type OrderedProduct struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    uint32  `json:"quantity"`
}
