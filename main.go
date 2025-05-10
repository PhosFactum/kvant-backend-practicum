package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/PhosFactum/kvant-backend-practicum/docs" // 👈 подключаем сгенерированные docs
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)


// Swagger API docs info

// @title Kvant API
// @version 1.0
// @description Task for a KVANT practicum
// @host localhost:8080
// @BasePath /


// Models

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Age      int    `json:"age"`
	Password string `json:"password"`
}

type Order struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	UserID    uint    `json:"user_id"`
	Product   string  `json:"product"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	CreatedAt string  `json:"created_at" gorm:"autoCreateTime"`
}


// Database conections

var db *gorm.DB


func initDB() {
	dsn := "host=localhost user=kvant_user password=11111111 dbname=kvant_db port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}
	db.AutoMigrate(&User{}, &Order{})
	fmt.Println("Connected to DB and migrated.")
}


// Handlers

// GetAllUsers godoc
// @Summary Get all users
// @Produce json
// @Success 200 {array} User
// @Failure 500 {object} map[string]string
// @Router /users [get]
func getUsers(c *gin.Context) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 404 {object} map[string]string
// @Router /user/{id} [get]
func getUserByID(c *gin.Context) {
	id := c.Param("id")
	log.Printf("Searching for user with ID: %s", id)
	var user User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// CreateUser godoc
// @Summary Create a new user
// @Accept json
// @Produce json
// @Param user body User true "User to create"
// @Success 201 {object} User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user [post]
func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var existing User
	if err := db.Where("email = ?", user.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	}
	if err := db.Create(&user).Error; err != nil {
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
// @Param user body User true "Updated user data"
// @Success 200 {object} User
// @Failure 404 {object} map[string]string
// @Router /user/{id} [put]
func updateUser(c *gin.Context) {
	id := c.Param("id")
	var user User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var updatedData User
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Name = updatedData.Name
	user.Email = updatedData.Email
	user.Age = updatedData.Age

	if err := db.Save(&user).Error; err != nil {
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
// @Router /user/{id} [delete]
func deleteUser(c *gin.Context) {
	id := c.Param("id")
	var user User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if err := db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.Status(http.StatusNoContent)
}


// Main entrypoint

func main() {
	initDB()
	router := gin.Default()

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/users", getUsers)
	router.GET("/user/:id", getUserByID)
	router.POST("/user", createUser)
	router.PUT("/user/:id", updateUser)
	router.DELETE("/user/:id", deleteUser)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}

