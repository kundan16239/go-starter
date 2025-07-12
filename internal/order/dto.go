package order

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateOrderRequest represents the request payload for creating an order
type CreateOrderRequest struct {
	UserID      string  `json:"user_id" validate:"required"`
	Items       []Item  `json:"items" validate:"required,min=1"`
	TotalAmount float64 `json:"total_amount" validate:"required,gt=0"`
}

// UpdateOrderRequest represents the request payload for updating an order
type UpdateOrderRequest struct {
	Status      string  `json:"status,omitempty" validate:"omitempty,oneof=pending processing shipped delivered cancelled"`
	TotalAmount float64 `json:"total_amount,omitempty" validate:"omitempty,gt=0"`
}

// OrderResponse represents the response payload for order data
type OrderResponse struct {
	ID          primitive.ObjectID `json:"id"`
	UserID      primitive.ObjectID `json:"user_id"`
	Items       []Item             `json:"items"`
	Status      string             `json:"status"`
	TotalAmount float64            `json:"total_amount"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// OrderListResponse represents the response payload for listing orders
type OrderListResponse struct {
	Orders []OrderResponse `json:"orders"`
	Total  int64           `json:"total"`
	Page   int             `json:"page"`
	Limit  int             `json:"limit"`
}

// Item represents an order item
type Item struct {
	ProductID   string  `json:"product_id" bson:"product_id"`
	ProductName string  `json:"product_name" bson:"product_name"`
	Quantity    int     `json:"quantity" bson:"quantity"`
	Price       float64 `json:"price" bson:"price"`
	Subtotal    float64 `json:"subtotal" bson:"subtotal"`
}

// OrderStatsResponse represents the response payload for order statistics
type OrderStatsResponse struct {
	TotalOrders  int64                    `json:"total_orders"`
	TotalRevenue float64                  `json:"total_revenue"`
	RecentOrders int64                    `json:"recent_orders"`
	StatusStats  []map[string]interface{} `json:"status_stats"`
	MonthlyStats []map[string]interface{} `json:"monthly_stats"`
}

// OrderAnalyticsResponse represents the response payload for order analytics
type OrderAnalyticsResponse struct {
	Days               int                      `json:"days"`
	DailyStats         []map[string]interface{} `json:"daily_stats"`
	TopCustomers       []map[string]interface{} `json:"top_customers"`
	AvgOrderValueStats []map[string]interface{} `json:"avg_order_value_stats"`
}

// ProductAnalyticsResponse represents the response payload for product analytics
type ProductAnalyticsResponse struct {
	TopProducts   []map[string]interface{} `json:"top_products"`
	AvgItemsStats []map[string]interface{} `json:"avg_items_stats"`
}
