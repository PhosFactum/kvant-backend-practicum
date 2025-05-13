package services

import (
    "context"
    "errors"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"

    "github.com/PhosFactum/kvant-backend-practicum/internal/models"
    "github.com/PhosFactum/kvant-backend-practicum/internal/repository"
)

var (
    ErrAuthInvalidCredentials = errors.New("invalid email or password")
)

// AuthService defines authentication use-cases.
type AuthService interface {
    Login(ctx context.Context, input models.LoginInput) (string, error)
}

// IsAuthError helps handlers identify auth errors.
func IsAuthError(err error) bool {
    return errors.Is(err, ErrAuthInvalidCredentials)
}

// authService is AuthService implementation.
type authService struct {
    userRepo repository.UserRepository
}

// NewAuthService constructs AuthService.
func NewAuthService(userRepo repository.UserRepository) AuthService {
    return &authService{userRepo: userRepo}
}

// Login implements password check and JWT creation.
func (s *authService) Login(ctx context.Context, input models.LoginInput) (string, error) {
    user, err := s.userRepo.FindByEmail(ctx, input.Email)
    if err != nil {
        return "", ErrAuthInvalidCredentials
    }
    if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)) != nil {
        return "", ErrAuthInvalidCredentials
    }

    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        secret = "supersecret"
    }

    claims := jwt.MapClaims{
        "user_id": user.ID,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signed, err := token.SignedString([]byte(secret))
    if err != nil {
        return "", err
    }
    return signed, nil
}

