package handlers

import (
    "net/http"
    "strconv"
    "log"

    "github.com/PhosFactum/kvant-backend-practicum/internal/models"
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    "golang.org/x/crypto/bcrypt"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
    db *gorm.DB
}

// NewUserHandler instantiates a new UserHandler
func NewUserHandler(db *gorm.DB) *UserHandler {
    return &UserHandler{db: db}
}

// GetUsers godoc
// @Summary Get all users with pagination and age filtering
// @Tags Users
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
    pageStr := c.DefaultQuery("page", "1")
    limitStr := c.DefaultQuery("limit", "10")
    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        log.Printf("GetUsers: invalid page parameter: %s", pageStr)
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page parameter"})
        return
    }
    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit < 1 {
        log.Printf("GetUsers: invalid limit parameter: %s", limitStr)
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
        return
    }

    // Parse age filters
    minAge := 0
    if minStr := c.Query("min_age"); minStr != "" {
        if v, err := strconv.Atoi(minStr); err != nil || v < 0 {
            log.Printf("GetUsers: invalid min_age parameter: %s", minStr)
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid min_age parameter"})
            return
        } else {
            minAge = v
        }
    }
    maxAge := 0
    if maxStr := c.Query("max_age"); maxStr != "" {
        if v, err := strconv.Atoi(maxStr); err != nil || v < 0 {
            log.Printf("GetUsers: invalid max_age parameter: %s", maxStr)
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid max_age parameter"})
            return
        } else {
            maxAge = v
        }
    }

    log.Printf("GetUsers: page=%d, limit=%d, min_age=%d, max_age=%d", page, limit, minAge, maxAge)

    query := h.db.Model(&models.User{})
    if minAge > 0 {
        query = query.Where("age >= ?", minAge)
    }
    if maxAge > 0 {
        query = query.Where("age <= ?", maxAge)
    }

    var total int
    if err := query.Count(&total).Error; err != nil {
        log.Printf("GetUsers: count error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count users"})
        return
    }

    offset := (page - 1) * limit
    var users []models.User
    if err := query.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
        log.Printf("GetUsers: find error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve users"})
        return
    }

    log.Printf("GetUsers: returned %d users", len(users))
    c.JSON(http.StatusOK, gin.H{
        "data":    users,
        "page":    page,
        "limit":   limit,
        "total":   total,
        "min_age": minAge,
        "max_age": maxAge,
    })
}

// GetUserByID godoc
// @Summary Get user by ID
// @Tags Users
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Router /user/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
    id := c.Param("id")
    log.Printf("GetUserByID: id=%s", id)
    var user models.User
    if err := h.db.First(&user, id).Error; err != nil {
        log.Printf("GetUserByID: not found: %s", id)
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }
    log.Printf("GetUserByID: found user id=%d", user.ID)
    c.JSON(http.StatusOK, user)
}

// CreateUser godoc
// @Summary Create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.CreateUserInput true "User to create"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    var input models.CreateUserInput
    if err := c.ShouldBindJSON(&input); err != nil {
        log.Printf("CreateUser: bind error: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    log.Printf("CreateUser: input email=%s, name=%s", input.Email, input.Name)

    var existing models.User
    if h.db.Where("email = ?", input.Email).First(&existing).RowsAffected > 0 {
        log.Printf("CreateUser: email exists: %s", input.Email)
        c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
        return
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        log.Printf("CreateUser: hash error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
        return
    }

    user := models.User{
        Name:         input.Name,
        Email:        input.Email,
        Age:          input.Age,
        PasswordHash: string(hash),
    }

    if err := h.db.Create(&user).Error; err != nil {
        log.Printf("CreateUser: create error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
        return
    }

    log.Printf("CreateUser: created user id=%d", user.ID)
    c.JSON(http.StatusCreated, gin.H{
        "id":    user.ID,
        "name":  user.Name,
        "email": user.Email,
        "age":   user.Age,
    })
}

// UpdateUser godoc
// @Summary Update user by ID
// @Tags Users
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param user body models.User true "Updated user data"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Router /user/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
    id := c.Param("id")
    log.Printf("UpdateUser: id=%s", id)
    var user models.User
    if err := h.db.First(&user, id).Error; err != nil {
        log.Printf("UpdateUser: not found: %s", id)
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    var updatedData models.User
    if err := c.ShouldBindJSON(&updatedData); err != nil {
        log.Printf("UpdateUser: bind error: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user.Name = updatedData.Name
    user.Email = updatedData.Email
    user.Age = updatedData.Age

    if err := h.db.Save(&user).Error; err != nil {
        log.Printf("UpdateUser: save error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
        return
    }

    log.Printf("UpdateUser: updated user id=%d", user.ID)
    c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete user by ID
// @Tags Users
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /user/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
    id := c.Param("id")
    log.Printf("DeleteUser: id=%s", id)
    var user models.User
    if err := h.db.First(&user, id).Error; err != nil {
        log.Printf("DeleteUser: not found: %s", id)
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    if err := h.db.Delete(&user).Error; err != nil {
        log.Printf("DeleteUser: delete error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
        return
    }

    log.Printf("DeleteUser: deleted user id=%d", user.ID)
    c.Status(http.StatusNoContent)
}

