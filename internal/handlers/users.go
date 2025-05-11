// internal/handlers/users.go
package handlers

import (
    "log"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    "golang.org/x/crypto/bcrypt"

    "github.com/PhosFactum/kvant-backend-practicum/internal/models"
)

// UserHandler handles user-related endpoints.
type UserHandler struct {
    db *gorm.DB
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(db *gorm.DB) *UserHandler {
    return &UserHandler{db: db}
}

// GetUsers retrieves users with pagination and optional age filters.
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
    page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
    if err != nil || page < 1 {
        log.Printf("GetUsers: invalid page: %q", c.Query("page"))
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
        return
    }
    limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
    if err != nil || limit < 1 {
        log.Printf("GetUsers: invalid limit: %q", c.Query("limit"))
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
        return
    }

    // Parse age filters
    minAge, maxAge := 0, 0
    if s := c.Query("min_age"); s != "" {
        if v, e := strconv.Atoi(s); e == nil && v >= 0 {
            minAge = v
        } else {
            log.Printf("GetUsers: invalid min_age: %q", s)
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid min_age"})
            return
        }
    }
    if s := c.Query("max_age"); s != "" {
        if v, e := strconv.Atoi(s); e == nil && v >= 0 {
            maxAge = v
        } else {
            log.Printf("GetUsers: invalid max_age: %q", s)
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid max_age"})
            return
        }
    }

    log.Printf("GetUsers: page=%d limit=%d min_age=%d max_age=%d", page, limit, minAge, maxAge)

    // Build query
    q := h.db.Model(&models.User{})
    if minAge > 0 {
        q = q.Where("age >= ?", minAge)
    }
    if maxAge > 0 {
        q = q.Where("age <= ?", maxAge)
    }

    // Count total
    var total int
    if err := q.Count(&total).Error; err != nil {
        log.Printf("GetUsers: count error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count users"})
        return
    }

    // Fetch page
    offset := (page - 1) * limit
    var users []models.User
    if err := q.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
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

// GetUserByID returns a single user by ID.
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

    log.Printf("GetUserByID: found user %d", user.ID)
    c.JSON(http.StatusOK, user)
}

// CreateUser registers a new user.
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
    var in models.CreateUserInput
    if err := c.ShouldBindJSON(&in); err != nil {
        log.Printf("CreateUser: bind error: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    log.Printf("CreateUser: email=%s name=%s", in.Email, in.Name)

    // Check duplicate email
    var exists models.User
    if h.db.Where("email = ?", in.Email).First(&exists).RowsAffected > 0 {
        log.Printf("CreateUser: email exists: %s", in.Email)
        c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
        return
    }

    // Hash password
    pw, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
    if err != nil {
        log.Printf("CreateUser: hash error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
        return
    }

    user := models.User{
        Name:         in.Name,
        Email:        in.Email,
        Age:          in.Age,
        PasswordHash: string(pw),
    }
    if err := h.db.Create(&user).Error; err != nil {
        log.Printf("CreateUser: db error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
        return
    }

    log.Printf("CreateUser: created id=%d", user.ID)
    c.JSON(http.StatusCreated, gin.H{
        "id":    user.ID,
        "name":  user.Name,
        "email": user.Email,
        "age":   user.Age,
    })
}

// UpdateUser modifies an existing user.
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
    id := c.Param("id")
    log.Printf("UpdateUser: id=%s", id)

    var user models.User
    if err := h.db.First(&user, id).Error; err != nil {
        log.Printf("UpdateUser: not found: %s", id)
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    var in models.User
    if err := c.ShouldBindJSON(&in); err != nil {
        log.Printf("UpdateUser: bind error: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user.Name = in.Name
    user.Email = in.Email
    user.Age = in.Age

    if err := h.db.Save(&user).Error; err != nil {
        log.Printf("UpdateUser: save error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
        return
    }

    log.Printf("UpdateUser: updated id=%d", user.ID)
    c.JSON(http.StatusOK, user)
}

// DeleteUser removes a user.
// @Summary Delete user by ID
// @Tags Users
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
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

    log.Printf("DeleteUser: deleted id=%d", user.ID)
    c.Status(http.StatusNoContent)
}

