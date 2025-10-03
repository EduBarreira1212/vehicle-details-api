package http

import (
	"github.com/EduBarreira1212/vehicle-details-api/internal/controllers"
	"github.com/gin-gonic/gin"
)

func BuildRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) { c.String(200, "Ok!") })

	api := r.Group("/api")

	api.POST("/users", controllers.CreateUser)
	api.GET("/users/:userID", controllers.GetUser)
	api.PUT("/users/:userID", controllers.UpdateUser)
	api.PUT("/users/:userID/update-password", controllers.UpdatePassword)
	api.DELETE("/users/:userID", controllers.DeleteUser)
	api.GET("/users/:userID/get-history", controllers.GetUserHistory)

	api.POST("/fipe/:userID", controllers.GetFipe)

	api.POST("/login", controllers.Login)

	return r
}
