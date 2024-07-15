package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/afroash/mastupeti/initializers"
	"github.com/afroash/mastupeti/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	// Get cookie of request
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"message": "Authorization token required"})
		return
	}

	// Decode and validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"message": "Invalid token"})
		return
	}

	// Extract claims and check token validity
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		// Check token expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatusJSON(401, gin.H{"message": "Token expired"})
			return
		}

		// Find user with token sub (subject)
		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatusJSON(402, gin.H{"message": "Invalid token"})
			return
		}

		// Attach user to request context
		c.Set("user", user)

		// Continue to the next handler
		c.Next()
	} else {
		c.AbortWithStatusJSON(403, gin.H{"message": "Invalid token"})
	}
}
