package http

import (
	"github.com/EduBarreira1212/vehicle-details-api/internal/controllers"
	"github.com/EduBarreira1212/vehicle-details-api/internal/models"
	"github.com/gin-gonic/gin"
)

func BuildRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) { c.String(200, "Ok!") })

	api := r.Group("/api")

	api.POST("/fipe", func(c *gin.Context) {
		var req models.FipeRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		c.JSON(200, gin.H{
			"fipe": req,
		})
	})

	api.POST("/users", controllers.CreateUser)
	api.GET("/users/:userID", controllers.GetUser)

	return r
}
