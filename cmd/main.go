package main

import (
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
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
    // Подключение к PostgreSQL
    db, err := gorm.Open("postgres", "postgres://kvant_user:11111111@localhost/kvant_db?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Автомиграции (опционально)
    db.AutoMigrate(&models.User{}, &models.Order{})

    // Инициализация Gin
    router := gin.Default()

    // Swagger
    docs.SwaggerInfo.BasePath = "/"
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Инициализация обработчиков
    userHandler := handlers.NewUserHandler(db)

    // Роуты
    router.GET("/users", userHandler.GetUsers)
    router.GET("/user/:id", userHandler.GetUserByID)
	router.POST("/user", userHandler.CreateUser)
    router.PUT("/user/:id", userHandler.UpdateUser)
    router.DELETE("/user/:id", userHandler.DeleteUser)

    // Запуск сервера
    if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatal(err)
    }
}
