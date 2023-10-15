package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	ProjectID      string
	AppEnv         string
	TopicID        string
	DeveloperEmail string
}

var Env *Environment

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	Env = &Environment{
		ProjectID:      getEnv("PROJECT_ID", "default-project-id"),
		AppEnv:         getEnv("APP_ENV", "development"),
		TopicID:        getEnv("TOPIC_ID", ""),
		DeveloperEmail: getEnv("DEVELOPER_EMAIL", ""),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
