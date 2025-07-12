package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	FirstName string `json:"first_name,omitempty" validate:"omitempty,min=2,max=50"`
	LastName  string `json:"last_name,omitempty" validate:"omitempty,min=2,max=50"`
	Email     string `json:"email,omitempty" validate:"omitempty,email"`
}

// UserResponse represents the response payload for user data
type UserResponse struct {
	ID        primitive.ObjectID `json:"id"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	Email     string             `json:"email"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// UserListResponse represents the response payload for listing users
type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

// LoginRequest represents the request payload for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the response payload for user login
type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// UserStatsResponse represents the response payload for user statistics
type UserStatsResponse struct {
	TotalUsers        int64                    `json:"total_users"`
	RecentUsers       int64                    `json:"recent_users"`
	RegistrationStats []map[string]interface{} `json:"registration_stats"`
	DomainStats       []map[string]interface{} `json:"domain_stats"`
}

// UserActivityResponse represents the response payload for user activity statistics
type UserActivityResponse struct {
	Days        int                      `json:"days"`
	DailyStats  []map[string]interface{} `json:"daily_stats"`
	HourlyStats []map[string]interface{} `json:"hourly_stats"`
}

// UserDemographicsResponse represents the response payload for user demographics
type UserDemographicsResponse struct {
	FirstNameLengthStats []map[string]interface{} `json:"first_name_length_stats"`
	LastNameLengthStats  []map[string]interface{} `json:"last_name_length_stats"`
}
