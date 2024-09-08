package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	//load environment variables
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
