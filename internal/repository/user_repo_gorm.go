package repository

import (
    "context"

    "github.com/jinzhu/gorm"
    "github.com/PhosFactum/kvant-backend-practicum/internal/models"
)

type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    FindByEmail(ctx context.Context, email string) (*models.User, error)
    List(ctx context.Context, page, limit, minAge, maxAge int) ([]models.User, int, error)
    GetByID(ctx context.Context, id uint) (*models.User, error)
    Update(ctx context.Context, user *models.User) error
    Delete(ctx context.Context, id uint) error
}

type gormUserRepo struct {
    db *gorm.DB
}

func NewGormUserRepo(db *gorm.DB) UserRepository {
    return &gormUserRepo{db: db}
}

func (r *gormUserRepo) Create(ctx context.Context, user *models.User) error {
    return r.db.Create(user).Error
}

func (r *gormUserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
    var user models.User
    if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *gormUserRepo) List(ctx context.Context, page, limit, minAge, maxAge int) ([]models.User, int, error) {
    var users []models.User
    var total int

    q := r.db.Model(&models.User{})
    if minAge > 0 {
        q = q.Where("age >= ?", minAge)
    }
    if maxAge > 0 {
        q = q.Where("age <= ?", maxAge)
    }

    if err := q.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    offset := (page - 1) * limit
    if err := q.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
        return nil, 0, err
    }

    return users, total, nil
}

func (r *gormUserRepo) GetByID(ctx context.Context, id uint) (*models.User, error) {
    var user models.User
    if err := r.db.First(&user, id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *gormUserRepo) Update(ctx context.Context, user *models.User) error {
    return r.db.Save(user).Error
}

func (r *gormUserRepo) Delete(ctx context.Context, id uint) error {
    return r.db.Delete(&models.User{}, id).Error
}

