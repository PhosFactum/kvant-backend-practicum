// models/auth.go
package models

// LoginInput represents the expected payload for user login
// swagger:model
// fields are bound and validated by Gin
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// TokenResponse defines the structure of the JWT token response
// swagger:model
type TokenResponse struct {
	Token string `json:"token"`
}

