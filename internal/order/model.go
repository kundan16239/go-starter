package order

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Order represents the order entity in the database
type Order struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Items       []Item             `json:"items" bson:"items"`
	Status      string             `json:"status" bson:"status"`
	TotalAmount float64            `json:"total_amount" bson:"total_amount"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// NewOrder creates a new order instance
func NewOrder(userID primitive.ObjectID, items []Item, totalAmount float64) *Order {
	return &Order{
		UserID:      userID,
		Items:       items,
		Status:      "pending",
		TotalAmount: totalAmount,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// UpdateStatus updates the order status
func (o *Order) UpdateStatus(status string) {
	o.Status = status
	o.UpdatedAt = time.Now()
}

// UpdateTotalAmount updates the order total amount
func (o *Order) UpdateTotalAmount(totalAmount float64) {
	o.TotalAmount = totalAmount
	o.UpdatedAt = time.Now()
}

// AddItem adds an item to the order
func (o *Order) AddItem(item Item) {
	o.Items = append(o.Items, item)
	o.UpdatedAt = time.Now()
}

// RemoveItem removes an item from the order by product ID
func (o *Order) RemoveItem(productID string) {
	for i, item := range o.Items {
		if item.ProductID == productID {
			o.Items = append(o.Items[:i], o.Items[i+1:]...)
			o.UpdatedAt = time.Now()
			break
		}
	}
}

// CalculateTotal recalculates the total amount based on items
func (o *Order) CalculateTotal() {
	var total float64
	for _, item := range o.Items {
		total += item.Subtotal
	}
	o.TotalAmount = total
	o.UpdatedAt = time.Now()
}

// IsValidStatus checks if the order status is valid
func (o *Order) IsValidStatus() bool {
	validStatuses := []string{"pending", "processing", "shipped", "delivered", "cancelled"}
	for _, status := range validStatuses {
		if o.Status == status {
			return true
		}
	}
	return false
}

// CanBeCancelled checks if the order can be cancelled
func (o *Order) CanBeCancelled() bool {
	return o.Status == "pending" || o.Status == "processing"
}

// CanBeShipped checks if the order can be shipped
func (o *Order) CanBeShipped() bool {
	return o.Status == "processing"
}

// CanBeDelivered checks if the order can be delivered
func (o *Order) CanBeDelivered() bool {
	return o.Status == "shipped"
}

// ToResponse converts the order to a response DTO
func (o *Order) ToResponse() OrderResponse {
	return OrderResponse{
		ID:          o.ID,
		UserID:      o.UserID,
		Items:       o.Items,
		Status:      o.Status,
		TotalAmount: o.TotalAmount,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}
}
