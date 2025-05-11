package handlers

import (
    "net/http"
    "github.com/PhosFactum/kvant-backend-practicum/internal/models"
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
)

type UserHandler struct {
    db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
    return &UserHandler{db: db}
}

// GetAllUsers godoc
// @Summary Get all users
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} map[string]string
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) { // Метод структуры UserHandler
    var users []models.User
    if err := h.db.Find(&users).Error; err != nil { // Используем h.db
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve users"})
        return
    }
    c.JSON(http.StatusOK, users)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Router /user/{id} [get] // Исправлен путь
func (h *UserHandler) GetUserByID(c *gin.Context) {
    id := c.Param("id")
    var user models.User
    if err := h.db.First(&user, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }
    c.JSON(http.StatusOK, user)
}

// CreateUser godoc
// @Summary Create a new user
// @Accept json
// @Produce json
// @Param user body models.User true "User to create"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post] // Исправлен путь
func (h *UserHandler) CreateUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Проверка уникальности email
    var existing models.User
    if h.db.Where("email = ?", user.Email).First(&existing).RowsAffected > 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
        return
    }

    if err := h.db.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "id":    user.ID,
        "name":  user.Name,
        "email": user.Email,
        "age":   user.Age,
    })
}

// UpdateUser godoc
// @Summary Update user by ID
// @Param id path int true "User ID"
// @Param user body models.User true "Updated user data"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Router /user/{id} [put] // Исправлен путь
func (h *UserHandler) UpdateUser(c *gin.Context) {
    id := c.Param("id")
    var user models.User
    if err := h.db.First(&user, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    var updatedData models.User
    if err := c.ShouldBindJSON(&updatedData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user.Name = updatedData.Name
    user.Email = updatedData.Email
    user.Age = updatedData.Age

    if err := h.db.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
        return
    }

    c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete user by ID
// @Param id path int true "User ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /user/{id} [delete] // Исправлен путь
func (h *UserHandler) DeleteUser(c *gin.Context) {
    id := c.Param("id")
    var user models.User
    if err := h.db.First(&user, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    if err := h.db.Delete(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
        return
    }

    c.Status(http.StatusNoContent)
}
