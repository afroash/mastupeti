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
	initializers.SyncUserDB()
}

func main() {
	r := gin.Default()
	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("templates/*")

	// Public routes
	r.GET("/", controllers.HomePage)

	//Signup and login routes
	r.POST("/signup", controllers.SignUp)
	r.GET("/signup", controllers.SignUp)
	r.POST("/login", controllers.LoginPage)
	r.GET("/login", controllers.LoginPage)

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
	}

	// Public video routes
	r.GET("/videos", controllers.VideoList)
	r.GET("/videos/:id", controllers.VideoShow)

	r.Run()
}
