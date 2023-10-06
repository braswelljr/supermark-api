package store

import (
	"time"
)

type User struct {
	ID          string    `json:"id" db:"id"`
	Fullname    string    `json:"fullname" db:"fullname"`
	Username    string    `json:"username" db:"username"`
	Email       string    `json:"email" db:"email"`
	DateOfBirth string    `json:"dateOfBirth" db:"date_of_birth"`
	Password    string    `json:"password" db:"password"`
	Phone       string    `json:"phone" db:"phone"`
	Roles       []string  `json:"roles" db:"roles"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type SignupPayload struct {
	Fullname    string `json:"fullname" validate:"required"`         // required
	Username    string `json:"username" validate:"required"`         // required
	DateOfBirth string `json:"dateOfBirth" validate:"omitempty"`     // optional
	Email       string `json:"email" validate:"required,email"`      // required
	Password    string `json:"password" validate:"required" min:"8"` // required
	Phone       string `json:"phone" validate:"required"`            // required
}

type UpdatePayload struct {
	Fullname    string `json:"fullname" db:"fullname" validate:"omitempty"`       // not required
	Username    string `json:"username" db:"username" validate:"omitempty"`       // not required
	DateOfBirth string `json:"dateOfBirth" db:"dateOfBirth" validate:"omitempty"` // not required
	Email       string `json:"email" db:"email" validate:"omitempty,email"`       // not required
	Phone       string `json:"phone" db:"phone" validate:"omitempty"`             // not required
}

type UpdateRolePayload struct {
	ID string `json:"id" db:"id" validate:"required"` // required
}

type UserResponse struct {
	ID          string    `json:"id" db:"id"`
	Fullname    string    `json:"fullname" db:"fullname"`
	Username    string    `json:"username" db:"username"`
	DateOfBirth string    `json:"dateOfBirth" db:"dateOfBirth"`
	Email       string    `json:"email" db:"email"`
	Phone       string    `json:"phone" db:"phone"`
	Roles       []string  `json:"roles" db:"roles"`
	CreatedAt   time.Time `json:"createdAt" db:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updatedAt"`
}

type PaginatedUsersResponse struct {
	Users       []UserResponse `json:"data"`
	Total       int            `json:"total" db:"total"`
	TotalPages  int            `json:"totalPages" db:"totalPages"`
	CurrentPage int            `json:"currentPage" db:"currentPage"`
}

type UserUpdateResponse struct {
	Message string `json:"message"`
}

type DeleteUserEvent struct {
	ID string `json:"id" db:"id"`
}
