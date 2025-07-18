package order

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateOrder(ctx context.Context, req CreateOrderRequest) (*OrderResponse, error) {
	order := &Order{
		Item:   req.Item,
		Amount: req.Amount,
	}
	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}
	resp := order.ToResponse()
	return &resp, nil
}

func (s *Service) GetOrderByID(ctx context.Context, id string) (*OrderResponse, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}
	order, err := s.repo.GetByID(ctx, oid)
	if err != nil {
		return nil, err
	}
	resp := order.ToResponse()
	return &resp, nil
}

func (s *Service) UpdateOrder(ctx context.Context, id string, req UpdateOrderRequest) (*OrderResponse, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}
	order, err := s.repo.GetByID(ctx, oid)
	if err != nil {
		return nil, err
	}
	if req.Item != "" {
		order.Item = req.Item
	}
	if req.Amount != 0 {
		order.Amount = req.Amount
	}
	if err := s.repo.Update(ctx, order); err != nil {
		return nil, err
	}
	resp := order.ToResponse()
	return &resp, nil
}

func (s *Service) DeleteOrder(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid order ID")
	}
	return s.repo.Delete(ctx, oid)
}

func (s *Service) ListOrders(ctx context.Context, skip, limit int64) (*OrderListResponse, error) {
	if limit < 1 || limit > 100 {
		limit = 10
	}
	findOptions := options.Find().SetSkip(skip).SetLimit(limit)
	orders, total, err := s.repo.List(ctx, skip, limit, findOptions)
	if err != nil {
		return nil, err
	}
	resp := make([]OrderResponse, len(orders))
	for i, o := range orders {
		resp[i] = o.ToResponse()
	}
	return &OrderListResponse{
		Orders: resp,
		Count:  total,
	}, nil
}
