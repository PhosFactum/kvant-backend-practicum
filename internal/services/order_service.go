package services

import (
    "context"
    "errors"
    "fmt"

    "github.com/PhosFactum/kvant-backend-practicum/internal/models"
    "github.com/PhosFactum/kvant-backend-practicum/internal/repository"
)

var (
    ErrUserNotFound   = errors.New("user not found")
    ErrInvalidRequest = errors.New("invalid request data")
)

// OrderService describes use-cases around orders.
type OrderService interface {
    Create(ctx context.Context, userID uint, req models.OrderRequest) (models.Order, error)
    ListByUser(ctx context.Context, userID uint) ([]models.Order, error)
    NotifyOrderCreated(ctx context.Context, order *models.Order) error
}

// IsNotFound helps handler map not-found errors.
func IsNotFound(err error) bool {
    return errors.Is(err, ErrUserNotFound)
}

// IsValidation helps handler map validation errors.
func IsValidation(err error) bool {
    return errors.Is(err, ErrInvalidRequest)
}

type orderService struct {
    userRepo  repository.UserRepository
    orderRepo repository.OrderRepository
}

// NewOrderService constructs OrderService.
func NewOrderService(u repository.UserRepository, o repository.OrderRepository) OrderService {
    return &orderService{userRepo: u, orderRepo: o}
}

func (s *orderService) Create(ctx context.Context, userID uint, req models.OrderRequest) (models.Order, error) {
    // 1) Проверка, что пользователь существует
    _, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return models.Order{}, ErrUserNotFound
    }

    // 2) Простейшая валидация полей
    if req.Quantity < 1 || req.Price < 0 {
        return models.Order{}, ErrInvalidRequest
    }

    // 3) Создание и сохранение
    order := models.Order{
        UserID:   userID,
        Product:  req.Product,
        Quantity: req.Quantity,
        Price:    req.Price,
    }
    if err := s.orderRepo.Create(ctx, &order); err != nil {
        return models.Order{}, fmt.Errorf("failed to create order: %w", err)
    }
    return order, nil
}

func (s *orderService) ListByUser(ctx context.Context, userID uint) ([]models.Order, error) {
    // 1) Убедиться, что пользователь есть
    _, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return nil, ErrUserNotFound
    }
    // 2) Вернуть заказы
    return s.orderRepo.ListByUser(ctx, userID)
}

// NotifyOrderCreated simulates sending a notification about the new order.
func (s *orderService) NotifyOrderCreated(ctx context.Context, order *models.Order) error {
    // Здесь можно интегрироваться с email/SMS/через сторонние сервисы
    fmt.Printf("Notifying user %d about new order %d\n", order.UserID, order.ID)
    return nil
}

