package controllers

import (
	"github.com/afroash/mastupeti/initializers"
	"github.com/afroash/mastupeti/models"
	"github.com/gin-gonic/gin"
)

func Dashboard(c *gin.Context) {

	type DashboardData struct {
		Username   string
		Email      string
		VideoCount int64
	}

	//get the current user
	user := c.MustGet("user").(models.User)

	var videoCount int64
	initializers.DB.Model(&models.Video{}).Where("user_id = ?", user.ID).Count(&videoCount)

	d := DashboardData{
		Username:   user.Username,
		Email:      user.Email,
		VideoCount: videoCount,
	}

	// render the dashboard
	c.HTML(200, "dashboard.html", d)
}
