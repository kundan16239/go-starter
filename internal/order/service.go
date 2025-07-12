package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"project-structure/pkg/shared/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service contains business logic for order operations
type Service struct {
	repo Repository
}

// NewService creates a new order service instance
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateOrder creates a new order with validation
func (s *Service) CreateOrder(ctx context.Context, req CreateOrderRequest) (*OrderResponse, error) {
	// Validate input
	if err := s.validateCreateOrderRequest(req); err != nil {
		return nil, err
	}

	// Convert user ID string to ObjectID
	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Create new order
	order := NewOrder(userID, req.Items, req.TotalAmount)

	// Save to database
	err = s.repo.Create(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	response := order.ToResponse()
	return &response, nil
}

// GetOrderByID retrieves an order by ID
func (s *Service) GetOrderByID(ctx context.Context, id string) (*OrderResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	order, err := s.repo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	response := order.ToResponse()
	return &response, nil
}

// UpdateOrder updates an existing order
func (s *Service) UpdateOrder(ctx context.Context, id string, req UpdateOrderRequest) (*OrderResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	// Validate input
	if err := s.validateUpdateOrderRequest(req); err != nil {
		return nil, err
	}

	// Get existing order
	order, err := s.repo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	// Update order fields
	if req.Status != "" {
		order.UpdateStatus(req.Status)
	}
	if req.TotalAmount > 0 {
		order.UpdateTotalAmount(req.TotalAmount)
	}

	// Save to database
	err = s.repo.Update(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	response := order.ToResponse()
	return &response, nil
}

// DeleteOrder removes an order from the system
func (s *Service) DeleteOrder(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid order ID")
	}

	return s.repo.Delete(ctx, objectID)
}

// ListOrders retrieves a paginated list of orders
func (s *Service) ListOrders(ctx context.Context, page, limit int) (*OrderListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	orders, total, err := s.repo.List(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	// Convert to response format
	orderResponses := make([]OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = order.ToResponse()
	}

	return &OrderListResponse{
		Orders: orderResponses,
		Total:  total,
		Page:   page,
		Limit:  limit,
	}, nil
}

// GetOrdersByUser retrieves orders for a specific user
func (s *Service) GetOrdersByUser(ctx context.Context, userID string, page, limit int) (*OrderListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	orders, total, err := s.repo.GetByUserID(ctx, objectID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}

	// Convert to response format
	orderResponses := make([]OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = order.ToResponse()
	}

	return &OrderListResponse{
		Orders: orderResponses,
		Total:  total,
		Page:   page,
		Limit:  limit,
	}, nil
}

// GetOrderStats retrieves order statistics using aggregation
func (s *Service) GetOrderStats(ctx context.Context) (*OrderStatsResponse, error) {
	// Get orders by status
	statusPipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":   "$status",
				"count": bson.M{"$sum": 1},
				"total": bson.M{"$sum": "$total_amount"},
			},
		},
		{
			"$sort": bson.M{"count": -1},
		},
	}

	statusStats, err := s.repo.Aggregate(ctx, statusPipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get status stats: %w", err)
	}

	// Get total order count and revenue
	totalOrders, err := s.repo.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get total orders: %w", err)
	}

	// Get total revenue
	revenuePipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":   nil,
				"total": bson.M{"$sum": "$total_amount"},
			},
		},
	}

	revenueStats, err := s.repo.Aggregate(ctx, revenuePipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue stats: %w", err)
	}

	var totalRevenue float64
	if len(revenueStats) > 0 {
		if total, ok := revenueStats[0]["total"].(float64); ok {
			totalRevenue = total
		}
	}

	// Get orders by month
	monthlyPipeline := []bson.M{
		{
			"$group": bson.M{
				"_id": bson.M{
					"year":  bson.M{"$year": "$created_at"},
					"month": bson.M{"$month": "$created_at"},
				},
				"count": bson.M{"$sum": 1},
				"total": bson.M{"$sum": "$total_amount"},
			},
		},
		{
			"$sort": bson.M{"_id.year": 1, "_id.month": 1},
		},
	}

	monthlyStats, err := s.repo.Aggregate(ctx, monthlyPipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly stats: %w", err)
	}

	// Get orders created in last 30 days
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	recentOrders, err := s.repo.CountDocuments(ctx, bson.M{
		"created_at": bson.M{"$gte": thirtyDaysAgo},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get recent orders: %w", err)
	}

	return &OrderStatsResponse{
		TotalOrders:  totalOrders,
		TotalRevenue: totalRevenue,
		RecentOrders: recentOrders,
		StatusStats:  bsonMToMapSlice(statusStats),
		MonthlyStats: bsonMToMapSlice(monthlyStats),
	}, nil
}

// GetOrderAnalytics retrieves detailed order analytics
func (s *Service) GetOrderAnalytics(ctx context.Context, days int) (*OrderAnalyticsResponse, error) {
	if days <= 0 {
		days = 30
	}

	startDate := time.Now().AddDate(0, 0, -days)

	// Get daily order counts and revenue
	dailyPipeline := []bson.M{
		{
			"$match": bson.M{
				"created_at": bson.M{"$gte": startDate},
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"$dateToString": bson.M{
						"format": "%Y-%m-%d",
						"date":   "$created_at",
					},
				},
				"count":         bson.M{"$sum": 1},
				"revenue":       bson.M{"$sum": "$total_amount"},
				"avgOrderValue": bson.M{"$avg": "$total_amount"},
			},
		},
		{
			"$sort": bson.M{"_id": 1},
		},
	}

	dailyStats, err := s.repo.Aggregate(ctx, dailyPipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily stats: %w", err)
	}

	// Get top customers by order count
	topCustomersPipeline := []bson.M{
		{
			"$match": bson.M{
				"created_at": bson.M{"$gte": startDate},
			},
		},
		{
			"$group": bson.M{
				"_id":           "$user_id",
				"orderCount":    bson.M{"$sum": 1},
				"totalSpent":    bson.M{"$sum": "$total_amount"},
				"avgOrderValue": bson.M{"$avg": "$total_amount"},
			},
		},
		{
			"$sort": bson.M{"totalSpent": -1},
		},
		{
			"$limit": 10,
		},
	}

	topCustomers, err := s.repo.Aggregate(ctx, topCustomersPipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get top customers: %w", err)
	}

	// Get average order value by status
	avgOrderValuePipeline := []bson.M{
		{
			"$match": bson.M{
				"created_at": bson.M{"$gte": startDate},
			},
		},
		{
			"$group": bson.M{
				"_id":      "$status",
				"avgValue": bson.M{"$avg": "$total_amount"},
				"minValue": bson.M{"$min": "$total_amount"},
				"maxValue": bson.M{"$max": "$total_amount"},
			},
		},
		{
			"$sort": bson.M{"avgValue": -1},
		},
	}

	avgOrderValueStats, err := s.repo.Aggregate(ctx, avgOrderValuePipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get average order value stats: %w", err)
	}

	return &OrderAnalyticsResponse{
		Days:               days,
		DailyStats:         bsonMToMapSlice(dailyStats),
		TopCustomers:       bsonMToMapSlice(topCustomers),
		AvgOrderValueStats: bsonMToMapSlice(avgOrderValueStats),
	}, nil
}

// GetProductAnalytics retrieves product-related analytics
func (s *Service) GetProductAnalytics(ctx context.Context) (*ProductAnalyticsResponse, error) {
	// Get top products by quantity sold
	topProductsPipeline := []bson.M{
		{
			"$unwind": "$items",
		},
		{
			"$group": bson.M{
				"_id":           "$items.product_id",
				"productName":   bson.M{"$first": "$items.product_name"},
				"totalQuantity": bson.M{"$sum": "$items.quantity"},
				"totalRevenue":  bson.M{"$sum": bson.M{"$multiply": []interface{}{"$items.price", "$items.quantity"}}},
				"orderCount":    bson.M{"$sum": 1},
			},
		},
		{
			"$sort": bson.M{"totalQuantity": -1},
		},
		{
			"$limit": 20,
		},
	}

	topProducts, err := s.repo.Aggregate(ctx, topProductsPipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get top products: %w", err)
	}

	// Get average items per order
	avgItemsPipeline := []bson.M{
		{
			"$addFields": bson.M{
				"itemCount": bson.M{"$size": "$items"},
			},
		},
		{
			"$group": bson.M{
				"_id":      nil,
				"avgItems": bson.M{"$avg": "$itemCount"},
				"minItems": bson.M{"$min": "$itemCount"},
				"maxItems": bson.M{"$max": "$itemCount"},
			},
		},
	}

	avgItemsStats, err := s.repo.Aggregate(ctx, avgItemsPipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get average items stats: %w", err)
	}

	return &ProductAnalyticsResponse{
		TopProducts:   bsonMToMapSlice(topProducts),
		AvgItemsStats: bsonMToMapSlice(avgItemsStats),
	}, nil
}

// validateCreateOrderRequest validates the create order request
func (s *Service) validateCreateOrderRequest(req CreateOrderRequest) error {
	if helpers.IsEmpty(req.UserID) {
		return errors.New("user ID is required")
	}

	if len(req.Items) == 0 {
		return errors.New("at least one item is required")
	}

	if req.TotalAmount <= 0 {
		return errors.New("total amount must be greater than 0")
	}

	// Validate each item
	for i, item := range req.Items {
		if helpers.IsEmpty(item.ProductID) {
			return fmt.Errorf("product ID is required for item %d", i+1)
		}
		if helpers.IsEmpty(item.ProductName) {
			return fmt.Errorf("product name is required for item %d", i+1)
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("quantity must be greater than 0 for item %d", i+1)
		}
		if item.Price <= 0 {
			return fmt.Errorf("price must be greater than 0 for item %d", i+1)
		}
	}

	return nil
}

// validateUpdateOrderRequest validates the update order request
func (s *Service) validateUpdateOrderRequest(req UpdateOrderRequest) error {
	if req.Status != "" {
		validStatuses := []string{"pending", "processing", "shipped", "delivered", "cancelled"}
		isValid := false
		for _, status := range validStatuses {
			if req.Status == status {
				isValid = true
				break
			}
		}
		if !isValid {
			return errors.New("invalid status")
		}
	}

	if req.TotalAmount < 0 {
		return errors.New("total amount cannot be negative")
	}

	return nil
}

func bsonMToMapSlice(input []bson.M) []map[string]interface{} {
	result := make([]map[string]interface{}, len(input))
	for i, v := range input {
		result[i] = v
	}
	return result
}
