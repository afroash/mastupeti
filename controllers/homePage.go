package controllers

import (
	"log"

	"github.com/afroash/mastupeti/initializers"
	"github.com/afroash/mastupeti/models"
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
	c.HTML(200, "index.html", gin.H{
		"videos": videos,
	})
}
