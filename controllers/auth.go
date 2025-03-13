package controllers

import (
	"fmt"
	"context"
	"net/http"
	"os"
	"restapis/config"
	"restapis/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Secret key for JWT (use environment variables in production)
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// RegisterUser - User Registration
func RegisterUser(c *gin.Context) {
	userCollection := config.GetCollection("users")
	if userCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not initialized"})
		return
	}

	var newUser models.User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("Received User:", newUser) // Debugging
	// Validation: Required fields
	if newUser.Email == "" || newUser.Password == "" || newUser.Role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email, password, and role are required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if email already exists
	var existingUser models.User
	err := userCollection.FindOne(ctx, bson.M{"email": newUser.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Hash the password
	hashedPassword, err := hashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	newUser.Password = hashedPassword
	newUser.ID = primitive.NewObjectID()

	_, err = userCollection.InsertOne(ctx, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	// Print success message in terminal
	fmt.Printf("User registered successfully: Email=%s, Role=%s\n", newUser.Email, newUser.Role)


	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login - Authenticate user and return JWT
func Login(c *gin.Context) {
	userCollection := config.GetCollection("users")
	if userCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not initialized"})
		return
	}

	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": loginData.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare hashed password
	if !checkPasswordHash(loginData.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := generateToken(user.ID.Hex(), user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	// Print success message in terminal
	fmt.Printf("User logged in successfully: Email=%s, Role=%s\n", user.Email, user.Role)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

// Generate JWT Token
func generateToken(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"role":   role,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// HashPassword - Encrypts user password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash - Verifies password hash
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
