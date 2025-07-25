package helpers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig(variableToReturn string) string {
	err := godotenv.Load("./../.env")
	if err != nil {
		log.Fatal("error loading the config file")
	}
	return os.Getenv(variableToReturn)
}
