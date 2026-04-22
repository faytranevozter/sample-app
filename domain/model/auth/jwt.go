package auth

import "github.com/golang-jwt/jwt/v5"

type JWTClaimUser struct {
	UserID string `json:"userID"`
	jwt.RegisteredClaims
}
