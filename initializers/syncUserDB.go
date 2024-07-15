package initializers

import "github.com/afroash/mastupeti/models"

func SyncUserDB() {
	DB.AutoMigrate(&models.User{})
}
