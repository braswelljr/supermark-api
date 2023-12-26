package store

import (
	"time"
)

type User struct {
	Id        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	Phone     string    `json:"phone" db:"phone"`
	Roles     []string  `json:"roles" db:"roles"`
	Avatar    string    `json:"avatar" db:"avatar"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type SignupPayload struct {
	Name     string `json:"name" validate:"required"`             // required
	Username string `json:"username" validate:"required"`         // required
	Email    string `json:"email" validate:"required,email"`      // required
	Phone    string `json:"phone" validate:"required"`            // required
	Password string `json:"password" validate:"required" min:"8"` // required
}

type UpdatePayload struct {
	Name     string `json:"name" db:"name" validate:"omitempty"`         // not required
	Username string `json:"username" db:"username" validate:"omitempty"` // not required
	Email    string `json:"email" db:"email" validate:"omitempty,email"` // not required
	Phone    string `json:"phone" db:"phone" validate:"omitempty"`       // not required
	Avatar   string `json:"avatar" db:"avatar" validate:"omitempty"`     // not required
}

type UpdateRolePayload struct {
	Id string `json:"id" db:"id" validate:"required"` // required
}

type UserResponse struct {
	Id        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Phone     string    `json:"phone" db:"phone"`
	Roles     []string  `json:"roles" db:"roles"`
	Avatar    string    `json:"avatar" db:"avatar"`
	CreatedAt time.Time `json:"createdAt" db:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" db:"updatedAt"`
}

type PaginatedUsersResponse struct {
	Users           []UserResponse `json:"data"`
	Total           int            `json:"total" db:"total"`
	TotalPages      int            `json:"totalPages" db:"totalPages"`
	CurrentPage     int            `json:"currentPage" db:"currentPage"`
	HasPreviousPage bool           `json:"hasPreviousPage" db:"hasPreviousPage"`
	HasNextPage     bool           `json:"hasNextPage" db:"hasNextPage"`
}

type UserUpdateResponse struct {
	Message string `json:"message"`
}

type DeleteUserEvent struct {
	Id string `json:"id" db:"id"`
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"` // required
	Password string `json:"password" validate:"required"`    // required
}

type Response struct {
	Message string        `json:"message"`
	Code    int           `json:"code"`
	Token   string        `json:"token"`
	Payload *UserResponse `json:"payload"`
}
