package http

import (
	"github.com/gin-gonic/gin"
)

func BuildRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) { c.String(200, "Ok!") })

	api := r.Group("/api")

	api.GET("/cars", func(c *gin.Context) { c.JSON(200, gin.H{"cars": []string{"car1", "car2", "car3"}}) })
	api.GET("/cars/:id", func(c *gin.Context) { c.JSON(200, gin.H{"car": 123}) })

	return r
}
