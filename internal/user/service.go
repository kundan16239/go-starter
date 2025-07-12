package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"project-structure/pkg/shared/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service contains business logic for user operations
type Service struct {
	repo Repository
}

// NewService creates a new user service instance
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateUser creates a new user with validation
func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
	// Validate input
	if err := s.validateCreateUserRequest(req); err != nil {
		return nil, err
	}
	// Check if user with email already exists
	exists, err := s.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, errors.New("user with this email already exists")
	}

	// Create new user
	user, err := NewUser(req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Save to database
	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

// GetUserByID retrieves a user by ID
func (s *Service) GetUserByID(ctx context.Context, id string) (*UserResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.repo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// UpdateUser updates an existing user
func (s *Service) UpdateUser(ctx context.Context, id string, req UpdateUserRequest) (*UserResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Validate input
	if err := s.validateUpdateUserRequest(req); err != nil {
		return nil, err
	}

	// Get existing user
	user, err := s.repo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	// Check if email is being updated and if it already exists
	if req.Email != "" && req.Email != user.Email {
		exists, err := s.repo.ExistsByEmail(ctx, req.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email existence: %w", err)
		}
		if exists {
			return nil, errors.New("user with this email already exists")
		}
	}

	// Update user fields
	user.Update(req.FirstName, req.LastName, req.Email)

	// Save to database
	err = s.repo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

// DeleteUser removes a user from the system
func (s *Service) DeleteUser(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user ID")
	}

	return s.repo.Delete(ctx, objectID)
}

// ListUsers retrieves a paginated list of users
func (s *Service) ListUsers(ctx context.Context, page, limit int) (*UserListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := s.repo.List(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Convert to response format
	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	return &UserListResponse{
		Users: userResponses,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

// LoginUser authenticates a user and returns user data
func (s *Service) LoginUser(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// Validate input
	if err := s.validateLoginRequest(req); err != nil {
		return nil, err
	}

	// Get user by email
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate token
	token, err := helpers.GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

// validateCreateUserRequest validates the create user request
func (s *Service) validateCreateUserRequest(req CreateUserRequest) error {
	if helpers.IsEmpty(req.FirstName) {
		return errors.New("first name is required")
	}
	if len(req.FirstName) < 2 || len(req.FirstName) > 50 {
		return errors.New("first name must be between 2 and 50 characters")
	}

	if helpers.IsEmpty(req.LastName) {
		return errors.New("last name is required")
	}
	if len(req.LastName) < 2 || len(req.LastName) > 50 {
		return errors.New("last name must be between 2 and 50 characters")
	}

	if helpers.IsEmpty(req.Email) {
		return errors.New("email is required")
	}
	if !helpers.IsValidEmail(req.Email) {
		return errors.New("invalid email format")
	}

	if helpers.IsEmpty(req.Password) {
		return errors.New("password is required")
	}
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	return nil
}

// validateUpdateUserRequest validates the update user request
func (s *Service) validateUpdateUserRequest(req UpdateUserRequest) error {
	if req.FirstName != "" && (len(req.FirstName) < 2 || len(req.FirstName) > 50) {
		return errors.New("first name must be between 2 and 50 characters")
	}

	if req.LastName != "" && (len(req.LastName) < 2 || len(req.LastName) > 50) {
		return errors.New("last name must be between 2 and 50 characters")
	}

	if req.Email != "" && !helpers.IsValidEmail(req.Email) {
		return errors.New("invalid email format")
	}

	return nil
}

// validateLoginRequest validates the login request
func (s *Service) validateLoginRequest(req LoginRequest) error {
	if helpers.IsEmpty(req.Email) {
		return errors.New("email is required")
	}
	if !helpers.IsValidEmail(req.Email) {
		return errors.New("invalid email format")
	}

	if helpers.IsEmpty(req.Password) {
		return errors.New("password is required")
	}

	return nil
}

// GetUserStats retrieves user statistics using aggregation
func (s *Service) GetUserStats(ctx context.Context) (*UserStatsResponse, error) {
	// Get user registration stats by month
	registrationPipeline := []bson.M{
		{
			"$group": bson.M{
				"_id": bson.M{
					"year":  bson.M{"$year": "$created_at"},
					"month": bson.M{"$month": "$created_at"},
				},
				"count": bson.M{"$sum": 1},
			},
		},
		{
			"$sort": bson.M{"_id.year": 1, "_id.month": 1},
		},
	}

	registrationStats, err := s.repo.Aggregate(ctx, registrationPipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get registration stats: %w", err)
	}

	// Get total user count
	totalUsers, err := s.repo.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get total users: %w", err)
	}

	// Get users created in last 30 days
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	recentUsers, err := s.repo.CountDocuments(ctx, bson.M{
		"created_at": bson.M{"$gte": thirtyDaysAgo},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get recent users: %w", err)
	}

	// Get users by email domain
	domainPipeline := []bson.M{
		{
			"$addFields": bson.M{
				"domain": bson.M{
					"$substr": []interface{}{
						"$email",
						bson.M{"$add": []interface{}{
							bson.M{"$indexOfBytes": []interface{}{"$email", "@"}},
							1,
						}},
						-1,
					},
				},
			},
		},
		{
			"$group": bson.M{
				"_id":   "$domain",
				"count": bson.M{"$sum": 1},
			},
		},
		{
			"$sort": bson.M{"count": -1},
		},
		{
			"$limit": 10,
		},
	}

	domainStats, err := s.repo.Aggregate(ctx, domainPipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain stats: %w", err)
	}

	return &UserStatsResponse{
		TotalUsers:        totalUsers,
		RecentUsers:       recentUsers,
		RegistrationStats: bsonMToMapSlice(registrationStats),
		DomainStats:       bsonMToMapSlice(domainStats),
	}, nil
}

// GetUserActivityStats retrieves user activity statistics
func (s *Service) GetUserActivityStats(ctx context.Context, days int) (*UserActivityResponse, error) {
	if days <= 0 {
		days = 30
	}

	startDate := time.Now().AddDate(0, 0, -days)

	// Get daily user registrations
	dailyRegistrationsPipeline := []bson.M{
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
				"count": bson.M{"$sum": 1},
			},
		},
		{
			"$sort": bson.M{"_id": 1},
		},
	}

	dailyStats, err := s.repo.Aggregate(ctx, dailyRegistrationsPipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily stats: %w", err)
	}

	// Get users by creation hour (to understand peak registration times)
	hourlyPipeline := []bson.M{
		{
			"$match": bson.M{
				"created_at": bson.M{"$gte": startDate},
			},
		},
		{
			"$group": bson.M{
				"_id":   bson.M{"$hour": "$created_at"},
				"count": bson.M{"$sum": 1},
			},
		},
		{
			"$sort": bson.M{"_id": 1},
		},
	}

	hourlyStats, err := s.repo.Aggregate(ctx, hourlyPipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get hourly stats: %w", err)
	}

	return &UserActivityResponse{
		Days:        days,
		DailyStats:  bsonMToMapSlice(dailyStats),
		HourlyStats: bsonMToMapSlice(hourlyStats),
	}, nil
}

// GetUserDemographics retrieves user demographic information
func (s *Service) GetUserDemographics(ctx context.Context) (*UserDemographicsResponse, error) {
	// Get users by first name length (simple demographic analysis)
	nameLengthPipeline := []bson.M{
		{
			"$addFields": bson.M{
				"firstNameLength": bson.M{"$strLenCP": "$first_name"},
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"$switch": bson.M{
						"branches": []bson.M{
							{"case": bson.M{"$lt": []interface{}{"$firstNameLength", 5}}, "then": "Short (1-4 chars)"},
							{"case": bson.M{"$lt": []interface{}{"$firstNameLength", 8}}, "then": "Medium (5-7 chars)"},
							{"case": bson.M{"$lt": []interface{}{"$firstNameLength", 12}}, "then": "Long (8-11 chars)"},
						},
						"default": "Very Long (12+ chars)",
					},
				},
				"count": bson.M{"$sum": 1},
			},
		},
		{
			"$sort": bson.M{"count": -1},
		},
	}

	nameLengthStats, err := s.repo.Aggregate(ctx, nameLengthPipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get name length stats: %w", err)
	}

	// Get users by last name length
	lastNameLengthPipeline := []bson.M{
		{
			"$addFields": bson.M{
				"lastNameLength": bson.M{"$strLenCP": "$last_name"},
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"$switch": bson.M{
						"branches": []bson.M{
							{"case": bson.M{"$lt": []interface{}{"$lastNameLength", 5}}, "then": "Short (1-4 chars)"},
							{"case": bson.M{"$lt": []interface{}{"$lastNameLength", 8}}, "then": "Medium (5-7 chars)"},
							{"case": bson.M{"$lt": []interface{}{"$lastNameLength", 12}}, "then": "Long (8-11 chars)"},
						},
						"default": "Very Long (12+ chars)",
					},
				},
				"count": bson.M{"$sum": 1},
			},
		},
		{
			"$sort": bson.M{"count": -1},
		},
	}

	lastNameLengthStats, err := s.repo.Aggregate(ctx, lastNameLengthPipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get last name length stats: %w", err)
	}

	return &UserDemographicsResponse{
		FirstNameLengthStats: bsonMToMapSlice(nameLengthStats),
		LastNameLengthStats:  bsonMToMapSlice(lastNameLengthStats),
	}, nil
}

// Helper to convert []bson.M to []map[string]interface{}
func bsonMToMapSlice(input []bson.M) []map[string]interface{} {
	result := make([]map[string]interface{}, len(input))
	for i, v := range input {
		result[i] = v
	}
	return result
}
