package db

import (
	"context"
	"fmt"
	"log"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func getSecret(secretID string) (string, error) {
	projectID := os.Getenv("ID_PROJECT")
	appEnv := os.Getenv("APP_ENV")

	log.Printf("Getting secret with ID: %s\n", secretID)
	log.Printf("Project ID: %s\n", projectID)
	log.Printf("App Environment: %s\n", appEnv)

	ctx := context.Background()

	var client *secretmanager.Client
	var err error

	if appEnv == "production" {
		client, err = secretmanager.NewClient(ctx)
		log.Println("Using production credentials for Secret Manager.")
	} else {
		client, err = secretmanager.NewClient(ctx, option.WithCredentialsFile("credentials.json"))
		log.Println("Using local credentials for Secret Manager.")
	}

	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretID)

	log.Printf("Constructed secret name: %s\n", name)

	result, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	})
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}

	return string(result.Payload.Data), nil
}

func DBConnection() {
	log.Println("Starting database connection...")

	host, err := getSecret("DB_HOST")
	if err != nil {
		log.Fatalf("Error getting DB_HOST: %v", err)
	}

	user, err := getSecret("DB_USER")
	if err != nil {
		log.Fatalf("Error getting DB_USER: %v", err)
	}

	password, err := getSecret("DB_PASSWORD")
	if err != nil {
		log.Fatalf("Error getting DB_PASSWORD: %v", err)
	}

	dbName, err := getSecret("DB_NAME")
	if err != nil {
		log.Fatalf("Error getting DB_NAME: %v", err)
	}

	port, err := getSecret("DB_PORT")
	if err != nil {
		log.Fatalf("Error getting DB_PORT: %v", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, password, dbName, port)

	log.Printf("Connecting to DB with DSN: %s\n", dsn) // Solo para debugging, no imprimas contraseñas en logs en entornos de producción.

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Database connected")
	}
}
