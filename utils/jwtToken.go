package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/afroash/mastupeti/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Utility function to generate which will return a  and validate JWT tokens

func GenerateToken(c *gin.Context, user models.User) (string, error) {
	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", fmt.Errorf("failed to upload file, %v", err)
	}

	return tokenString, nil
}

func ValidateToken(c *gin.Context, tokenString string) error {
	return nil
}
