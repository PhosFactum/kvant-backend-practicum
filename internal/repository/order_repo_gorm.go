package repository

import (
    "context"

    "github.com/jinzhu/gorm"
    "github.com/PhosFactum/kvant-backend-practicum/internal/models"
)

// OrderRepository defines DB operations for orders.
type OrderRepository interface {
    Create(ctx context.Context, order *models.Order) error
    ListByUser(ctx context.Context, userID uint) ([]models.Order, error)
}

type gormOrderRepo struct {
    db *gorm.DB
}

// NewGormOrderRepo creates a GORM implementation.
func NewGormOrderRepo(db *gorm.DB) OrderRepository {
    return &gormOrderRepo{db: db}
}

func (r *gormOrderRepo) Create(ctx context.Context, order *models.Order) error {
    return r.db.Create(order).Error
}

func (r *gormOrderRepo) ListByUser(ctx context.Context, userID uint) ([]models.Order, error) {
    var orders []models.Order
    if err := r.db.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
        return nil, err
    }
    return orders, nil
}

