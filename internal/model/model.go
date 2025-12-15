package model

import "github.com/golang-jwt/jwt/v4"

type Caller struct {
	Login string
}

type AuthRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Claims struct {
	jwt.RegisteredClaims
	Login string `json:"login"`
}
