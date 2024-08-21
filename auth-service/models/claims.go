package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	ID   string `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}
