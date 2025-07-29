package helpers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig(variableToReturn string) string {
	// Try current directory first, then parent directory
	err := godotenv.Load(".env")
	if err != nil {
		err = godotenv.Load("./../.env")
		if err != nil {
			log.Fatal("error loading the config file")
		}
	}
	return os.Getenv(variableToReturn)
}
