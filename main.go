package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Age      int    `json:"age"`
	Password string `json:"password"` // храним как есть для MVP
}

var db *gorm.DB

func initDB() {
	dsn := "host=localhost user=kvant_user password=11111111 dbname=kvant_db port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}
	db.AutoMigrate(&User{})
	fmt.Println("Connected to DB and migrated.")
}

func main() {
	initDB()

	router := gin.Default()

	router.POST("/user", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// проверка на уникальность email
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
	})

	router.GET("/user", func(c *gin.Context) {
		var users []User
		if err := db.Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve users"})
			return
		}
		c.JSON(http.StatusOK, users)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}

