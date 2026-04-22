package middleware

import (
	"app/domain/model/auth"
	"app/domain/model/response"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (m *appMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		hAuth := c.GetHeader("Authorization")
		if hAuth == "" {
			response := response.Error(http.StatusUnauthorized, "Unauthorized: Header authorization is required")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		splitToken := strings.Split(hAuth, "Bearer ")
		if len(splitToken) != 2 {
			response := response.Error(http.StatusUnauthorized, "Unauthorized: Token is invalid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// get token without 'Bearer '
		tokenString := splitToken[1]

		// validating token
		token, err := jwt.ParseWithClaims(tokenString, &auth.JWTClaimUser{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.secret), nil
		})

		// check validity token
		if token == nil || !token.Valid {
			if errors.Is(err, jwt.ErrTokenMalformed) {
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					response.Error(http.StatusUnauthorized, "Unauthorized: Token is invalid"),
				)
				return
			}

			if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					response.Error(http.StatusUnauthorized, "Unauthorized: Token signature invalid"),
				)
				return
			}

			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					response.Error(http.StatusUnauthorized, "Unauthorized: Token expired"),
				)
				return
			}

			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				response.Error(http.StatusUnauthorized, err.Error()),
			)
			return
		}

		claims, tokenOK := token.Claims.(*auth.JWTClaimUser)
		if !tokenOK {
			response := response.Error(http.StatusUnauthorized, "Unauthorized: Token data not valid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("token_data", *claims)
		c.Next()
	}
}
