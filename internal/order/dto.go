package order

type CreateOrderRequest struct {
	Item   string  `json:"item" validate:"required"`
	Amount float64 `json:"amount" validate:"required"`
}

type UpdateOrderRequest struct {
	Item   string  `json:"item" validate:"omitempty"`
	Amount float64 `json:"amount" validate:"omitempty"`
}

type OrderResponse struct {
	ID     string  `json:"id"`
	Item   string  `json:"item"`
	Amount float64 `json:"amount"`
}

type OrderListResponse struct {
	Orders []OrderResponse `json:"orders"`
	Count  int64           `json:"count"`
}
