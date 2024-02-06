package initialisers

import (
	"github.com/web-auth-go/03_jwt/models"
	"log"
)

func SyncDatabase() {
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("could not run migrate user")
	}
}
