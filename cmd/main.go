package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"

	"github.com/PhosFactum/kvant-backend-practicum/docs"
	"github.com/PhosFactum/kvant-backend-practicum/internal/handlers"
	"github.com/PhosFactum/kvant-backend-practicum/internal/middleware"
	"github.com/PhosFactum/kvant-backend-practicum/internal/models"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Kvant API
// @version 1.0
// @description Task for a KVANT practicum
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @host localhost:8080
// @BasePath /
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка при загрузке .env файла")
	}

	// Setup database connection
	db := initDB()
	defer db.Close()

	// AutoMigrate models
	db.AutoMigrate(&models.User{}, &models.Order{})

	// Init Gin router
	router := gin.Default()

	// Swagger
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Handlers
	userHandler := handlers.NewUserHandler(db)
	orderHandler := handlers.NewOrderHandler(db)
	authHandler := handlers.NewAuthHandler(db)

	// Public routes
	router.POST("/auth/login", authHandler.Login)
	router.POST("/users", userHandler.CreateUser)
	router.GET("/users", userHandler.GetUsers)
	router.GET("/user/:id", userHandler.GetUserByID)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.JWTAuthMiddleware())

	protected.PUT("/user/:id", userHandler.UpdateUser)
	protected.DELETE("/user/:id", userHandler.DeleteUser)

	userGroup := protected.Group("/users/:user_id")
	{
		userGroup.POST("/orders", orderHandler.CreateOrder)
		userGroup.GET("/orders", orderHandler.GetOrdersByUser)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}

// initDB sets up the PostgreSQL connection
func initDB() *gorm.DB {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable TimeZone=UTC",
		dbHost, dbPort, dbUser, dbName, dbPassword)

	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	return db
}

