package main

import (
	"github.com/afroash/mastupeti/initializers"
	"github.com/afroash/mastupeti/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
}
func main() {
	initializers.DB.AutoMigrate(&models.Video{})
}
