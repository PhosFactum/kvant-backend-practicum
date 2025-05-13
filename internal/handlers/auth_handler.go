package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/PhosFactum/kvant-backend-practicum/internal/models"
    "github.com/PhosFactum/kvant-backend-practicum/internal/services"
)

// AuthHandler manages authentication-related endpoints.
type AuthHandler struct {
    svc services.AuthService
}

// NewAuthHandler returns a new AuthHandler.
func NewAuthHandler(svc services.AuthService) *AuthHandler {
    return &AuthHandler{svc: svc}
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

    token, err := h.svc.Login(c.Request.Context(), input)
    switch {
    case services.IsAuthError(err):
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
    case err != nil:
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    default:
        c.JSON(http.StatusOK, models.TokenResponse{Token: token})
    }
}

