// internal/handlers/user_handler.go
package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/PhosFactum/kvant-backend-practicum/internal/models"
    "github.com/PhosFactum/kvant-backend-practicum/internal/services"
)

type UserHandler struct {
    svc services.UserService
}

func NewUserHandler(svc services.UserService) *UserHandler {
    return &UserHandler{svc: svc}
}

// @Summary Get all users with pagination and age filtering
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(10)
// @Param min_age query int false "Minimum age to filter"
// @Param max_age query int false "Maximum age to filter"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    minAge, _ := strconv.Atoi(c.DefaultQuery("min_age", "0"))
    maxAge, _ := strconv.Atoi(c.DefaultQuery("max_age", "0"))

    users, total, err := h.svc.List(c.Request.Context(), page, limit, minAge, maxAge)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data":    users,
        "total":   total,
        "page":    page,
        "limit":   limit,
        "min_age": minAge,
        "max_age": maxAge,
    })
}

// @Summary Get user by ID
// @Tags Users
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Router /user/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    user, err := h.svc.GetByID(c.Request.Context(), uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, user)
}

// @Summary Create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.CreateUserInput true "User to create"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    var input models.CreateUserInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.svc.Create(c.Request.Context(), input)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, user)
}

// @Summary Update user by ID
// @Tags Users
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param user body models.User true "Updated user data"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    var input models.User
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    updated, err := h.svc.Update(c.Request.Context(), uint(id), input)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, updated)
}

// @Summary Delete user by ID
// @Tags Users
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}

