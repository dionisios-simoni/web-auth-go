package initialisers

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	DB, err = gorm.Open(postgres.Open(os.Getenv("DB")), &gorm.Config{})
	if err != nil {
		panic("failed to connect to db")
	}

}
