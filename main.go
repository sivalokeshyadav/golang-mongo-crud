package main

import (
	"log"
	"restapis/config"
	"restapis/controllers"
	"restapis/middleware"
	"restapis/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to Database
	if err := config.ConnectDB(); err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Create a new Gin router
	router := gin.Default()

	// Public Routes
	router.POST("/register", controllers.RegisterUser)
	router.POST("/login", controllers.Login)

	// Protected Routes (Authenticated Users Only)
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/users", controllers.GetUsers)

	// Admin Routes (Require "admin" Role)
	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
	admin.DELETE("/user/:id", controllers.DeleteUser)

	// Initialize Additional Routes
	routes.UserRoutes(router)

	// Start Server
	log.Println("🚀 Server running on port 8080")
	router.Run(":8080")
}
