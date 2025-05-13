package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    _ "github.com/lib/pq"

    "github.com/PhosFactum/kvant-backend-practicum/docs"
    "github.com/PhosFactum/kvant-backend-practicum/internal/handlers"
    "github.com/PhosFactum/kvant-backend-practicum/internal/middleware"
    "github.com/PhosFactum/kvant-backend-practicum/internal/models"
    "github.com/PhosFactum/kvant-backend-practicum/internal/repository"
    "github.com/PhosFactum/kvant-backend-practicum/internal/services"

    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

// @title KVANT Backend Practicum API
// @version 1.0
// @description REST API для управления пользователями и их заказами с JWT-авторизацией

// @contact.name KVANT Team
// @contact.url https://kvant-is.ru/

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите JWT-токен в формате: `Bearer <token>`

// @tag.name Users
// @tag.description Операции с пользователями

// @tag.name Orders
// @tag.description Управление заказами пользователей

// @tag.name Auth
// @tag.description Аутентификация и получение JWT-токена

// @x-logo {"url": "https://kvant-team.com/logo.png", "backgroundColor": "#FFFFFF", "altText": "KVANT Logo"}
func main() {
    // Initialize DB connection
    db := initDB()
    defer db.Close()

    // Migrate schema
    db.AutoMigrate(&models.User{}, &models.Order{})

    // Initialize repositories
    userRepo := repository.NewGormUserRepo(db)
    orderRepo := repository.NewGormOrderRepo(db)

    // Initialize services
    authSvc := services.NewAuthService(userRepo)
    userSvc := services.NewUserService(userRepo)
    orderSvc := services.NewOrderService(userRepo, orderRepo)

    // Initialize handlers
    authH := handlers.NewAuthHandler(authSvc)
    userH := handlers.NewUserHandler(userSvc)
    orderH := handlers.NewOrderHandler(orderSvc)

    // Setup router
    router := gin.Default()
    docs.SwaggerInfo.BasePath = "/"
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Public endpoints
    router.POST("/auth/login", authH.Login)
    router.POST("/users", userH.CreateUser)
    router.GET("/users", userH.GetUsers)
    router.GET("/user/:id", userH.GetUserByID)

    // Protected endpoints
    protected := router.Group("/")
    protected.Use(middleware.JWTAuthMiddleware())
    {
        protected.PUT("/user/:id", userH.UpdateUser)
        protected.DELETE("/user/:id", userH.DeleteUser)

        userGroup := protected.Group("/users/:user_id")
        {
            userGroup.POST("/orders", orderH.CreateOrder)
            userGroup.GET("/orders", orderH.GetOrdersByUser)
        }
    }

    // Start HTTP server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    addr := fmt.Sprintf(":%s", port)
    log.Printf("starting server on %s", addr)
    if err := http.ListenAndServe(addr, router); err != nil {
        log.Fatal("server failed:", err)
    }
}

// initDB builds the Postgres DSN and opens a GORM connection
func initDB() *gorm.DB {
    host := os.Getenv("DB_HOST")
    user := os.Getenv("DB_USER")
    pass := os.Getenv("DB_PASSWORD")
    name := os.Getenv("DB_NAME")
    port := os.Getenv("DB_PORT")

    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s dbname=%s password=%s sslmode=disable TimeZone=UTC",
        host, port, user, name, pass,
    )

    db, err := gorm.Open("postgres", dsn)
    if err != nil {
        log.Fatal("database connection failed:", err)
    }
    return db
}

