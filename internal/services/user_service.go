package services

import (
    "context"
    "errors"
    "fmt"

    "golang.org/x/crypto/bcrypt"

    "github.com/PhosFactum/kvant-backend-practicum/internal/models"
    "github.com/PhosFactum/kvant-backend-practicum/internal/repository"
)

type UserService interface {
    Create(ctx context.Context, input models.CreateUserInput) (*models.User, error)
    List(ctx context.Context, page, limit, minAge, maxAge int) ([]models.User, int, error)
    GetByID(ctx context.Context, id uint) (*models.User, error)
    Update(ctx context.Context, id uint, input models.User) (*models.User, error)
    Delete(ctx context.Context, id uint) error
    SendWelcomeEmail(ctx context.Context, user *models.User) error
}

type userService struct {
    repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
    return &userService{repo: r}
}

func (s *userService) Create(ctx context.Context, input models.CreateUserInput) (*models.User, error) {
    exists, _ := s.repo.FindByEmail(ctx, input.Email)
    if exists != nil {
        return nil, errors.New("email already exists")
    }

    pwHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    user := &models.User{
        Name:         input.Name,
        Email:        input.Email,
        Age:          input.Age,
        PasswordHash: string(pwHash),
    }

    err = s.repo.Create(ctx, user)
    return user, err
}

func (s *userService) List(ctx context.Context, page, limit, minAge, maxAge int) ([]models.User, int, error) {
    return s.repo.List(ctx, page, limit, minAge, maxAge)
}

func (s *userService) GetByID(ctx context.Context, id uint) (*models.User, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *userService) Update(ctx context.Context, id uint, input models.User) (*models.User, error) {
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    user.Name = input.Name
    user.Email = input.Email
    user.Age = input.Age

    err = s.repo.Update(ctx, user)
    return user, err
}

func (s *userService) Delete(ctx context.Context, id uint) error {
    return s.repo.Delete(ctx, id)
}

// SendWelcomeEmail simulates sending a welcome email to the new user.
func (s *userService) SendWelcomeEmail(ctx context.Context, user *models.User) error {
    // Здесь может быть интеграция с email-сервисом.
    // Пока просто логируем.
    fmt.Printf("Sending welcome email to %s at %s\n", user.Name, user.Email)
    return nil
}

