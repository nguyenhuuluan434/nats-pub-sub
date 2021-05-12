package model

import "encoding/json"

type Order struct {
	OrderID    int
	CustomerID string
	Status     string
}

func NewOrder(orderId int, customerId string, status string) *Order {
	return &Order{
		OrderID:    orderId,
		CustomerID: customerId,
		Status:     status,
	}
}
func (order *Order) ToJson() ([]byte, error) {
	return json.Marshal(order)
}
