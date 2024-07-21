package main

import (
	"github.com/afroash/mastupeti/controllers"
	"github.com/afroash/mastupeti/initializers"
	"github.com/afroash/mastupeti/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
	//initializers.SyncUserDB()
}

func main() {
	r := gin.Default()

	r.Static("/assests", "./assests")
	r.LoadHTMLGlob("templates/*")

	// Public routes
	r.GET("/", controllers.HomePage)
	//to be done
	r.GET("/about", controllers.AboutPage)
	r.GET("/contact", controllers.ContactPage)
	r.POST("/submitform", controllers.ContactForm)

	r.GET("/admin", controllers.AdminPage)

	//Signup and login routes
	r.POST("/signup", controllers.SignUp)
	r.GET("/signup", controllers.SignUp)
	r.POST("/login", controllers.LoginPage)
	r.GET("/login", controllers.LoginPage)

	//Modal routes
	r.GET("/signup-modal", controllers.SignUpModal)
	r.GET("/login-modal", controllers.LoginModal)
	r.GET("/upload-modal", controllers.UploadFormModal)

	// Routes that require authentication
	auth := r.Group("/")
	auth.Use(middleware.RequireAuth)
	{
		auth.GET("/validate", controllers.Validate)
		auth.POST("/logout", controllers.Logout)
		auth.GET("/logout", controllers.Logout)
		auth.GET("/addvideo", controllers.VideoCreate)
		auth.POST("/videos", controllers.VideoCreate)
		auth.PUT("/videos/:id", controllers.VideoUpdate)
		auth.DELETE("/videos/:id", controllers.VideoDelete)
		auth.GET("/stats", controllers.Dashboard)
	}

	// Public video routes
	r.GET("/videos", controllers.VideoList)
	r.GET("/videos/:id", controllers.VideoShow)

	r.Run()
}
