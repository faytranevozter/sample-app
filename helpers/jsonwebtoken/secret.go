package jwt_helper

import (
	"app/domain/model/auth"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GetJwtCredential() JWTCredential {
	ttl, _ := strconv.Atoi(os.Getenv("JWT_MEMBER_TTL"))

	return JWTCredential{
		Member: JWT{
			Secret: os.Getenv("JWT_MEMBER_SECRET_KEY"),
			TLL:    time.Duration(ttl) * time.Minute,
			Algo:   jwt.SigningMethodHS256,
		},
	}
}

func GenerateJWTToken(jwtCred JWT, data jwt.Claims) (string, error) {
	var newClaims jwt.Claims

	// reassign data
	if claim, ok := data.(auth.JWTClaimUser); ok {
		if claim.RegisteredClaims.ID == "" {
			claim.RegisteredClaims.ID = uuid.NewString()
		}

		if claim.RegisteredClaims.Issuer == "" {
			claim.RegisteredClaims.Issuer = "member"
		}

		if claim.RegisteredClaims.IssuedAt == nil {
			claim.RegisteredClaims.IssuedAt = jwt.NewNumericDate(time.Now())
		}

		if claim.RegisteredClaims.NotBefore == nil {
			claim.RegisteredClaims.NotBefore = jwt.NewNumericDate(time.Now())
		}

		if claim.RegisteredClaims.ExpiresAt == nil && jwtCred.TLL.String() != "0s" {
			claim.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(jwtCred.TLL))
		}

		newClaims = claim
	} else {
		return "", fmt.Errorf("claim data not supported")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(jwtCred.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
