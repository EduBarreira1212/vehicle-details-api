package middlewares

import (
	"net/http"

	"github.com/EduBarreira1212/vehicle-details-api/internal/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := auth.ValidateToken(c); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		if uid, err := auth.ExtractUserID(c); err == nil {
			c.Set("userId", uid)
		}

		c.Next()
	}
}
