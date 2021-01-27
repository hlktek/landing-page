package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// ListConfig is list all dotenv variable
var ListConfig = map[string]string{}

func init() {
	e := godotenv.Load() //Load .env file
	if e != nil {
		log.Fatal(e)
	}
	ListConfig = GetAllConfig()
}

// GetConfig is get one
func GetConfig(key string) string {
	//GetConfig get config
	return os.Getenv(key)
}

// GetAllConfig is get all variable in dotenv
func GetAllConfig() map[string]string {
	listConfig, err := godotenv.Read()
	if err != nil {
		log.Fatal(err)
	}
	return listConfig
}
