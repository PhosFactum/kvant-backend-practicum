package models


type User struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	Name         string `json:"name"`
	Email        string `json:"email" gorm:"unique"`
	Age          int    `json:"age"`
	PasswordHash string `json:"-" gorm:"column:password_hash"`
}

type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Age      int    `json:"age" binding:"required"`
	Password string `json:"password" binding:"required"`
}

