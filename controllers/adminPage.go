package controllers

import (
	"log"
	"net/http"

	"github.com/afroash/mastupeti/initializers"
	"github.com/afroash/mastupeti/models"
	"github.com/gin-gonic/gin"
)

func AdminPage(c *gin.Context) {

	//Get All Users
	var users []models.User
	u := initializers.DB.Find(&users)
	if u.Error != nil {
		log.Printf("Error fetching users: %v", u.Error)
	}

	//Get All Videos
	var videos []models.Video
	v := initializers.DB.Find(&videos)
	if v.Error != nil {
		log.Printf("Error fetching videos: %v", v.Error)
	}

	//return all users and videos to the admin page
	c.HTML(http.StatusOK, "superadmin.html", gin.H{
		"users":  users,
		"videos": videos,
	})
}
