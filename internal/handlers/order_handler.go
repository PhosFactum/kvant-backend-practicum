// internal/handlers/order_handler.go
package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/PhosFactum/kvant-backend-practicum/internal/models"
    "github.com/PhosFactum/kvant-backend-practicum/internal/services"
    "github.com/PhosFactum/kvant-backend-practicum/internal/utils"
)

// OrderHandler manages order-related endpoints.
type OrderHandler struct {
    svc services.OrderService
}

// NewOrderHandler creates a new OrderHandler.
func NewOrderHandler(svc services.OrderService) *OrderHandler {
    return &OrderHandler{svc: svc}
}

// CreateOrder creates a new order for a given user.
// @Summary Create order for a user
// @Tags Orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param order body models.OrderRequest true "Order info"
// @Success 201 {object} models.Order
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{user_id}/orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
    userID, err := strconv.Atoi(c.Param("user_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
        return
    }

    var req models.OrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    order, err := h.svc.Create(c.Request.Context(), uint(userID), req)
    switch {
    case services.IsNotFound(err):
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    case services.IsValidation(err):
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    case err != nil:
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    default:
        // Асинхронное уведомление о новом заказе
        utils.Async(func() {
            _ = h.svc.NotifyOrderCreated(c.Request.Context(), &order)
        })
        c.JSON(http.StatusCreated, order)
    }
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
    userID, err := strconv.Atoi(c.Param("user_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
        return
    }

    orders, err := h.svc.ListByUser(c.Request.Context(), uint(userID))
    switch {
    case services.IsNotFound(err):
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    case err != nil:
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    default:
        c.JSON(http.StatusOK, orders)
    }
}

