package controllers

import (
	"net/http"

	"github.com/afroash/mastupeti/initializers"
	"github.com/afroash/mastupeti/models"
	"github.com/afroash/mastupeti/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//User controller packege.

func SignUp(c *gin.Context) {
	// Handle user signup (POST only)
	if c.Request.Method == "POST" {
		var body struct {
			Email    string `form:"email" binding:"required,email"`
			Password string `form:"password" binding:"required"`
			Username string `form:"username" binding:"required"`
		}

		if err := c.ShouldBind(&body); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		//Check if the email is already in the database.
		var existingUser models.User
		initializers.DB.Where("email = ?", body.Email).First(&existingUser)
		if existingUser.ID != 0 {
			c.JSON(400, gin.H{"error": "Email already in use"})
			return
		}

		// Check if the username is already in the database.
		initializers.DB.Where("username = ?", body.Username).First(&existingUser)
		if existingUser.ID != 0 {
			c.JSON(400, gin.H{"error": "Username already in use"})
			return
		}

		// Hash the password (you'll need to import a hashing library)
		hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10) // Use a cost parameter of 10 or higher
		if err != nil {
			c.JSON(500, gin.H{"error": "Error hashing password"})
			return
		}

		// Create a new user
		user := models.User{Email: body.Email, Password: string(hash), Username: body.Username}
		result := initializers.DB.Create(&user)

		if result.Error != nil {
			c.JSON(500, gin.H{"error": "Failed to create user"})
			return
		}

		// Optionally, generate a JWT token for the user and send it in the response or store it in a cookie.

		// Redirect to login page or dashboard or show success message (your choice)
		c.HTML(200, "signUp.html", gin.H{
			"showform": false, // Pass success message to template
		})
	} else { // GET
		c.HTML(http.StatusOK, "signUp.html", gin.H{
			"showForm": true, // Indicate to template to show the form
		})
	}
}

// Login function
func LoginPage(c *gin.Context) {
	// Get the email and password from the form request.
	var body struct {
		Email    string `form:"email" binding:"required,email"`
		Password string `form:"password" binding:"required"`
	}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login data!!"})
		return
	}

	// Check if the email is in the database.
	var user models.User
	result := initializers.DB.First(&user, "email = ?", body.Email)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid email or password"}) // Generic message to avoid revealing if email exists
		return

	}

	// Check if the password is correct.
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	//if password is incorrect return login modal with error message
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return

	}

	// Generate a JWT token.
	tokenString, err := utils.GenerateToken(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Save JWT token in a cookie.
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600, "", "", false, true)

	// Redirect on successful login
	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func Validate(c *gin.Context) {
	_, exists := c.Get("user") // Check if user is logged in
	if exists {
		//user is logged in, show "logout" button
		c.HTML(http.StatusOK, "logoutButton.html", nil)
	} else {
		//user is not logged in, show "login" button
		c.HTML(http.StatusOK, "loginButton.html", nil)
	}

}

func Logout(c *gin.Context) {
	// Extract the token from the request's Authorization cookie
	token, err := c.Cookie("Authorization")
	if err == nil && token != "" {
		// Invalidate the token on the server-side
		utils.InvalidateToken(token)
	}

	// Set the Authorization cookie to expire immediately
	c.SetCookie("Authorization", "", -1, "", "", false, true)

	// After logout redirect to Home page
	c.Redirect(http.StatusFound, "/")

	//c.HTML(http.StatusOK, "loginButton.html", nil)
}
