package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    "github.com/PhosFactum/kvant-backend-practicum/internal/models"
)


type OrderRequest struct {
    Product  string  `json:"product" binding:"required" example:"Laptop"`
    Quantity int     `json:"quantity" binding:"required,min=1" example:"1"`
    Price    float64 `json:"price" binding:"required,min=0" example:"1200.50"`
}


type OrderHandler struct {
    db *gorm.DB
}


func NewOrderHandler(db *gorm.DB) *OrderHandler {
    return &OrderHandler{db: db}
}


// CreateOrder godoc
// @Summary Create order for a user
// @Description Creating a new order for a user
// @Tags Orders
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param order body OrderRequest true "Order info"
// @Success 201 {object} models.Order
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{user_id}/orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
    userID := c.Param("user_id")

    var user models.User
    if err := h.db.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    var orderReq OrderRequest
    if err := c.ShouldBindJSON(&orderReq); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    order := models.Order{
        UserID:   user.ID,
        Product:  orderReq.Product,
        Quantity: orderReq.Quantity,
        Price:    orderReq.Price,
    }

    if err := h.db.Create(&order).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
        return
    }

    c.JSON(http.StatusCreated, order)
}


// GetOrdersByUser godoc
// @Summary Get orders for a user
// @Description Checking all orders of a user
// @Tags Orders
// @Param user_id path int true "User ID"
// @Success 200 {array} models.Order
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{user_id}/orders [get]
func (h *OrderHandler) GetOrdersByUser(c *gin.Context) {
    userID := c.Param("user_id")

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

    if err := h.db.
        Where("user_id = ?", userID).
        Find(&orders).
        Error; err != nil {
    	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve orders"})
        return
    }

    c.JSON(http.StatusOK, orders)
}
