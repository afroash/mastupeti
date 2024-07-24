package controllers

import (
	"log"
	"net/http"

	"github.com/afroash/mastupeti/initializers"
	"github.com/afroash/mastupeti/middleware"
	"github.com/afroash/mastupeti/models"
	"github.com/afroash/mastupeti/utils"
	"github.com/gin-gonic/gin"
)

// function to be used to render the home page.
func HomePage(c *gin.Context) {
	//get the last 3 videos
	var videos []models.Video
	vid := initializers.DB.Order("created_at desc").Limit(3).Find(&videos)
	if vid.Error != nil {
		log.Printf("Error fetching videos: %v", vid.Error)
	}
	//check if there is a cookie with the name authToken

	isAuthenticated := middleware.IsAuthenticated(c)
	if isAuthenticated {

		//update the videos to include the signed URL
		keys := make([]string, len(videos))
		for i, video := range videos {
			keys[i] = video.URL
		}
		urls, err := utils.GetSignedURLs(keys)
		if err != nil {
			log.Printf("Error getting signed URLs: %v", err)
		}
		for i, video := range videos {
			videos[i].URL = urls[video.URL]
		}

		c.HTML(200, "index.html", gin.H{
			"videos":          videos,
			"IsAuthenticated": isAuthenticated,
		})
		return
	}

}

// function to render the signup modal.
func SignUpModal(c *gin.Context) {
	c.HTML(200, "signUpModal.html", nil)
}

// function to render the login modal.
func LoginModal(c *gin.Context) {
	c.HTML(200, "loginModal.html", nil)
}

// function to render the login modal.
func UploadFormModal(c *gin.Context) {
	c.HTML(200, "uploadVideoModal.html", nil)
}

// function to render about modal.
func AboutPage(c *gin.Context) {
	c.HTML(200, "about.html", nil)
}

// function to render the contact modal.
func ContactPage(c *gin.Context) {
	c.HTML(200, "contact.html", nil)
}

// func to submit the contact form
func ContactForm(c *gin.Context) {
	var body struct {
		Email   string `form:"email" binding:"required,email"`
		Subject string `form:"subject" binding:"required"`
		Message string `form:"message" binding:"required"`
	}
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	//send the email
	if err := utils.SendEmail(body.Email, body.Subject, body.Message, "ashleyc@masterash.co.uk"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}
	c.Redirect(302, "/")
}
