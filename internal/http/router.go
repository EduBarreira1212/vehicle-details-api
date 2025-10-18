package http

import (
	"github.com/EduBarreira1212/vehicle-details-api/internal/controllers"
	"github.com/EduBarreira1212/vehicle-details-api/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func BuildRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) { c.String(200, "Ok!") })

	api := r.Group("/api")

	api.GET("/me", middlewares.AuthMiddleware(), controllers.GetMyProfile)
	api.POST("/users", controllers.CreateUser)
	api.GET("/users/:userID", middlewares.AuthMiddleware(), controllers.GetUser)
	api.PUT("/users/:userID", middlewares.AuthMiddleware(), controllers.UpdateUser)
	api.PUT("/users/:userID/update-password", middlewares.AuthMiddleware(), controllers.UpdatePassword)
	api.DELETE("/users/:userID", middlewares.AuthMiddleware(), controllers.DeleteUser)
	api.GET("/users/:userID/get-history", middlewares.AuthMiddleware(), controllers.GetUserHistory)

	api.POST("/fipe", middlewares.AuthMiddleware(), controllers.GetFipe)

	api.POST("/login", controllers.Login)

	return r
}
