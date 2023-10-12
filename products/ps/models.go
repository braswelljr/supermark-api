package ps

import "time"

type Product struct {
	Id          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price"`
	CategoryId  string    `json:"categoryId" db:"category_id"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type ProductRequest struct {
	Name          string  `json:"name" validate:"required"`
	Description   string  `json:"description"  validate:"required"`
	Price         float64 `json:"price" validate:"required,min=0"`
	CategoryId    string  `json:"categoryId"  validate:"omitempty"`
	StockQuantity int     `json:"stockQuantity" validate:"required,min=0"`
}
