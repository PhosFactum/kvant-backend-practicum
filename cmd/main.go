package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    "github.com/joho/godotenv"
    "github.com/PhosFactum/kvant-backend-practicum/internal/handlers"
    "github.com/PhosFactum/kvant-backend-practicum/internal/models"
    "github.com/PhosFactum/kvant-backend-practicum/docs"
    "github.com/jinzhu/gorm"
    _ "github.com/lib/pq"
)

// @title Kvant API
// @version 1.0
// @description Task for a KVANT practicum
// @host localhost:8080
// @BasePath /
func main() {
    // Загружаем переменные из .env
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Ошибка при загрузке .env файла")
    }

    dbHost := os.Getenv("DB_HOST")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    dbPort := os.Getenv("DB_PORT")
    port := os.Getenv("PORT")

    // Формируем строку подключения
    dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbName, dbPassword)

    // Подключение к PostgreSQL
    db, err := gorm.Open("postgres", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Автомиграции
    db.AutoMigrate(&models.User{}, &models.Order{})

    // Инициализация Gin
    router := gin.Default()

    // Swagger
    docs.SwaggerInfo.BasePath = "/"
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Инициализация обработчиков
    userHandler := handlers.NewUserHandler(db)
	orderHandler := handlers.NewOrderHandler(db)

    // Роуты
    router.GET("/users", userHandler.GetUsers)
    router.GET("/user/:id", userHandler.GetUserByID)
    router.POST("/users", userHandler.CreateUser)
    router.PUT("/user/:id", userHandler.UpdateUser)
    router.DELETE("/user/:id", userHandler.DeleteUser)

	userGroup := router.Group("/users/:user_id")
	{
		userGroup.POST("/orders", orderHandler.CreateOrder)
		userGroup.GET("/orders", orderHandler.GetOrdersByUser)
	}

    // Запуск сервера
    addr := fmt.Sprintf(":%s", port)
    if err := http.ListenAndServe(addr, router); err != nil {
        log.Fatal(err)
    }
}

