package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    "github.com/PhosFactum/kvant-backend-practicum/internal/models"
)

type OrderHandler struct {
    db *gorm.DB
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
    return &OrderHandler{db: db}
}

// CreateOrder godoc
// @Summary Create order for a user
// @Param user_id path int true "User ID"
// @Param order body models.Order true "Order data"
// @Success 201 {object} models.Order
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{user_id}/orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
    userID := c.Param("user_id")

    // Проверка, существует ли пользователь
    var user models.User
    if err := h.db.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    // Получение данных заказа
    var order models.Order
    if err := c.ShouldBindJSON(&order); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    order.UserID = user.ID
    if err := h.db.Create(&order).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
        return
    }

    c.JSON(http.StatusCreated, order)
}

// GetOrdersByUser godoc
// @Summary Get orders for a user
// @Param user_id path int true "User ID"
// @Success 200 {array} models.Order
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{user_id}/orders [get]
func (h *OrderHandler) GetOrdersByUser(c *gin.Context) {
    userID := c.Param("user_id")

    // Проверка, существует ли пользователь
    var user models.User
    if err := h.db.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    var orders []models.Order
    if err := h.db.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve orders"})
        return
    }

    c.JSON(http.StatusOK, orders)
}
