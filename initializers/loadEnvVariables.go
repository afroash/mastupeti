package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnvVariables loads the .env file
func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
