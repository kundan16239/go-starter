package user

import (
	"time"

	"project-structure/pkg/shared/helpers"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents the user entity in the database
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName string             `bson:"first_name" json:"first_name"`
	LastName  string             `bson:"last_name" json:"last_name"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// NewUser creates a new user instance
func NewUser(firstName, lastName, email, password string) (*User, error) {
	hashedPassword, err := helpers.HashPassword(password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// CheckPassword verifies if the provided password matches the user's password
func (u *User) CheckPassword(password string) bool {
	return helpers.CheckPasswordHash(password, u.Password)
}

// Update updates the user fields
func (u *User) Update(firstName, lastName, email string) {
	if firstName != "" {
		u.FirstName = firstName
	}
	if lastName != "" {
		u.LastName = lastName
	}
	if email != "" {
		u.Email = email
	}
	u.UpdatedAt = time.Now()
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
