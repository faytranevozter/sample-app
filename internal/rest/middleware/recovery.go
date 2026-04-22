package middleware

import (
	"app/domain/model/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (c *appMiddleware) Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.Error("Panic Recover : ", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Something went wrong"))
			}
		}()
		c.Next()
	}
}
