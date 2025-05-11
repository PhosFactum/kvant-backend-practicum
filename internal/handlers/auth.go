// internal/handlers/auth.go
package handlers

import (
    "net/http"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "github.com/jinzhu/gorm"
    "golang.org/x/crypto/bcrypt"

    "github.com/PhosFactum/kvant-backend-practicum/internal/models"
)

// AuthHandler manages authentication-related endpoints.
type AuthHandler struct {
    db *gorm.DB
}

// NewAuthHandler returns a new AuthHandler.
func NewAuthHandler(db *gorm.DB) *AuthHandler {
    return &AuthHandler{db: db}
}

// Login authenticates user credentials and returns a JWT.
// @Summary Login and get JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginInput true "Login credentials"
// @Success 200 {object} models.TokenResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
    var input models.LoginInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    if err := h.db.Where("email = ?", input.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
        return
    }

    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        secret = "supersecret"
    }

    // Create token with 24h expiry
    claims := jwt.MapClaims{
        "user_id": user.ID,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    signed, err := token.SignedString([]byte(secret))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign token"})
        return
    }

    c.JSON(http.StatusOK, models.TokenResponse{Token: signed})
}

