package order

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Item   string             `bson:"item" json:"item"`
	Amount float64            `bson:"amount" json:"amount"`
}

func (o *Order) ToResponse() OrderResponse {
	return OrderResponse{
		ID:     o.ID.Hex(),
		Item:   o.Item,
		Amount: o.Amount,
	}
}
