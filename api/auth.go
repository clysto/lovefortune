package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var ACCESS_KEY = os.Getenv("ACCESS_KEY")

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-Access-Token")
		if key == ACCESS_KEY {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"description": "没有权限访问",
			})
		}
	}
}
