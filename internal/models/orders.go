// models/orders.go
package models

import "time"

// Order represents a purchase order linked to a user
// swagger:model
type Order struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"-" gorm:"index"`
	Product   string    `json:"product"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP"`
}

// OrderRequest defines the payload for creating an order
// swagger:model
type OrderRequest struct {
	Product  string  `json:"product" binding:"required"`
	Quantity int     `json:"quantity" binding:"required,min=1"`
	Price    float64 `json:"price" binding:"required,min=0"`
}

