package jwt_helper

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	Secret string
	Issuer string
	TLL    time.Duration
	Algo   jwt.SigningMethod
}

type JWTCredential struct {
	Member JWT
}
