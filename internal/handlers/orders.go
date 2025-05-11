// internal/handlers/orders.go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    "github.com/PhosFactum/kvant-backend-practicum/internal/models"
)

// OrderHandler manages order-related endpoints.
type OrderHandler struct {
    db *gorm.DB
}

// NewOrderHandler creates a new OrderHandler.
func NewOrderHandler(db *gorm.DB) *OrderHandler {
    return &OrderHandler{db: db}
}

// CreateOrderRequest defines the input for creating an order.
type CreateOrderRequest struct {
    Product  string  `json:"product" binding:"required" example:"Laptop"`
    Quantity int     `json:"quantity" binding:"required,min=1" example:"1"`
    Price    float64 `json:"price" binding:"required,min=0" example:"1200.50"`
}

// CreateOrder creates a new order for a given user.
// @Summary Create order for a user
// @Tags Orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param order body CreateOrderRequest true "Order info"
// @Success 201 {object} models.Order
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{user_id}/orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
    userID := c.Param("user_id")

    // Ensure user exists
    var user models.User
    if err := h.db.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    // Bind request
    var req CreateOrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Save order
    order := models.Order{
        UserID:   user.ID,
        Product:  req.Product,
        Quantity: req.Quantity,
        Price:    req.Price,
    }
    if err := h.db.Create(&order).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
        return
    }

    c.JSON(http.StatusCreated, order)
}

// GetOrdersByUser returns all orders for a given user.
// @Summary Get orders for a user
// @Tags Orders
// @Security BearerAuth
// @Param user_id path int true "User ID"
// @Success 200 {array} models.Order
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{user_id}/orders [get]
func (h *OrderHandler) GetOrdersByUser(c *gin.Context) {
    userID := c.Param("user_id")

    // Ensure user exists
    var user models.User
    if err := h.db.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    // Fetch orders
    var orders []models.Order
    if err := h.db.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve orders"})
        return
    }

    c.JSON(http.StatusOK, orders)
}

