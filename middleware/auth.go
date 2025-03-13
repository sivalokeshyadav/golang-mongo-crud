package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"restapis/utils"
)

// // AuthMiddleware checks JWT token validity
// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		tokenString := c.GetHeader("Authorization")
// 		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 			c.Abort()
// 			return
// 		}

// 		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
// 		token, err := utils.ValidateToken(tokenString)
// 		if err != nil || !token.Valid {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 			c.Abort()
// 			return
// 		}

// 		// Store token claims in context for later use
// 		claims, _ := token.Claims.(jwt.MapClaims)
// 		c.Set("email", claims["email"])
// 		c.Set("role", claims["role"])
// 		c.Next()
// 	}
// }

// // RoleMiddleware ensures only authorized roles can access certain routes
// func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		role, exists := c.Get("role")
// 		if !exists {
// 			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
// 			c.Abort()
// 			return
// 		}

// 		for _, allowed := range allowedRoles {
// 			if role == allowed {
// 				c.Next()
// 				return
// 			}
// 		}

// 		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
// 		c.Abort()
// 	}
// }
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		// Log the received token
		println("Received Authorization Header:", tokenString)

		if !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, err := utils.ValidateToken(tokenString)
		if err != nil {
			println("Token validation error:", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store token claims in context for later use
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Set("email", claims["email"])
		c.Set("role", claims["role"])
		c.Next()
	}
}
// RoleMiddleware ensures only authorized roles can access certain routes
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		for _, allowed := range allowedRoles {
			if role == allowed {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}