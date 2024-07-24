package controllers

import (
	"log"

	"github.com/afroash/mastupeti/initializers"
	"github.com/afroash/mastupeti/middleware"
	"github.com/afroash/mastupeti/models"
	"github.com/afroash/mastupeti/utils"
	"github.com/gin-gonic/gin"
)

func VideoCreate(c *gin.Context) {
	// Handle video upload (POST only)
	var fileURL string
	if c.Request.Method == "POST" {
		file, fileHeader, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{"error": "File upload failed!"})
			return
		}
		defer file.Close()

		// Upload video to S3
		fileURL, err = utils.UploadToS3(file, fileHeader)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error uploading file to S3"})
			return
		}
	}
	// Get authenticated user from context
	user, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(500, gin.H{"message": "Internal server error"})
		return
	}

	userID := user.(models.User).ID // Extract the UserID from the User object

	// Get data from request (for both GET and POST)
	var body struct {
		Title string `json:"title" form:"title"`
		Body  string `json:"body" form:"body"`
	}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid form data"})
		return
	}

	// Handle both GET and POST
	if c.Request.Method == "POST" {
		// Create video (POST only)
		video := models.Video{
			Title:  body.Title,
			URL:    fileURL,
			Body:   body.Body,
			UserID: userID,
		}
		result := initializers.DB.Create(&video)
		if result.Error != nil {
			c.JSON(500, gin.H{"error": "Failed to create video"})
			return
		}

		// Return "Add Video" button HTML snippet to replace the form
		c.Redirect(302, "/videos")
	} else { // GET
		// Return the initial form for adding videos
		c.HTML(200, "addVideo.html", gin.H{
			"showForm": true, // Indicate to template to show the form
		})
	}
}
func VideoList(c *gin.Context) {
	// get videos
	var videos []models.Video
	vid := initializers.DB.Find(&videos)
	if vid.Error != nil {
		log.Printf("Error fetching videos: %v", vid.Error)
	}
	// return all videos
	// Render the index.html template with the videos data
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
func VideoShow(c *gin.Context) {
	// get video
	var video models.Video
	initializers.DB.First(&video, c.Param("id"))

	// return video
	c.JSON(200, gin.H{
		"video": video,
	})
}

func VideoUpdate(c *gin.Context) {
	// get id from url

	// get data from request
	var body struct {
		Title string
		URL   string
		Body  string
	}

	c.Bind(&body)

	// get video
	var video models.Video
	initializers.DB.First(&video, c.Param("id"))

	// update video
	video.Title = body.Title
	video.URL = body.URL
	video.Body = body.Body

	initializers.DB.Model(&video).Updates(models.Video{
		Title: body.Title,
		URL:   body.URL,
		Body:  body.Body,
	})

	// return video
	c.JSON(200, gin.H{
		"video": video,
	})
}

func VideoDelete(c *gin.Context) {
	// get id from url

	// get video
	var video models.Video
	initializers.DB.First(&video, c.Param("id"))

	// delete video
	initializers.DB.Delete(&video)

	// respond
	c.Status(200)
}
