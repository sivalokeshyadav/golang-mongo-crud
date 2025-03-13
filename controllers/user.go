package controllers

import (
	"context"
	"fmt"
	"net/http"
	"restapis/config"
	"restapis/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetUsers - Fetch all users (Protected Route)
func GetUsers(c *gin.Context) {
	userCollection := config.GetCollection("users")
	if userCollection == nil {
		fmt.Println("Error: Database not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not initialized"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []models.User
	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("Error: Failed to fetch users", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			fmt.Println("Error: Decoding user data", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding user data"})
			return
		}
		users = append(users, user)
	}

	fmt.Println("Success: Users fetched successfully")
	c.JSON(http.StatusOK, users)
}

// GetUser - Fetch user by ID
func GetUser(c *gin.Context) {
	userCollection := config.GetCollection("users")
	if userCollection == nil {
		fmt.Println("Error: Database not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not initialized"})
		return
	}

	userID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		fmt.Println("Error: Invalid user ID", userID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		fmt.Println("Error: User not found", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	fmt.Println("Success: User fetched successfully", user)
	c.JSON(http.StatusOK, user)
}

// UpdateUser - Update a user by ID
func UpdateUser(c *gin.Context) {
	userCollection := config.GetCollection("users")
	if userCollection == nil {
		fmt.Println("Error: Database not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not initialized"})
		return
	}

	userID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		fmt.Println("Error: Invalid user ID", userID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var updatedUser models.User
	if err := c.BindJSON(&updatedUser); err != nil {
		fmt.Println("Error: Failed to parse request body", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"name": updatedUser.Name, "email": updatedUser.Email}}
	_, err = userCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		fmt.Println("Error: Failed to update user", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	fmt.Println("Success: User updated successfully", updatedUser)
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser - Admin can delete user
func DeleteUser(c *gin.Context) {
	userCollection := config.GetCollection("users")
	if userCollection == nil {
		fmt.Println("Error: Database not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not initialized"})
		return
	}

	userID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		fmt.Println("Error: Invalid user ID", userID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = userCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		fmt.Println("Error: Failed to delete user", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	fmt.Println("Success: User deleted successfully", userID)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
