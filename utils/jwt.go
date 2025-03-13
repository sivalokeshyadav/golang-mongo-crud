// package utils

// import (
// 	"errors"
// 	"time"

// 	"github.com/golang-jwt/jwt/v5"
// )

// var jwtSecret = []byte("your_secret_key") // Change this to a secure secret

// // GenerateToken creates a JWT token
// func GenerateToken(email, role string) (string, error) {
// 	claims := jwt.MapClaims{
// 		"email": email,
// 		"role":  role,
// 		"exp":   time.Now().Add(time.Hour * 2).Unix(), // Token expires in 2 hours
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString(jwtSecret)
// }

// // ValidateToken verifies the JWT token
// func ValidateToken(tokenString string) (*jwt.Token, error) {
// 	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, errors.New("invalid token")
// 		}
// 		return jwtSecret, nil
// 	})
// }
package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"os"
	"fmt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(email, role string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	fmt.Printf("Generating token with claims: %+v\n", claims) // Debugging

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		fmt.Println("Error signing token:", err)
	}
	return signedToken, err
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		fmt.Println("Token validation error:", err)
	} else {
		fmt.Printf("Decoded Token: %+v\n", token.Claims) // Debugging
	}

	return token, err
}
