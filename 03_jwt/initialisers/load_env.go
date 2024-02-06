package initialisers

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("unable to load environment variables from .env")
	}
}
