package routes

import (
	"restapis/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/", controllers.RegisterUser) // Changed from CreateUser to RegisterUser
		userGroup.GET("/", controllers.GetUsers)
		userGroup.GET("/:id", controllers.GetUser)
		userGroup.PUT("/:id", controllers.UpdateUser) // Ensure UpdateUser is defined in user.go
		userGroup.DELETE("/:id", controllers.DeleteUser)
	}
}
