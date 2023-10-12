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
	// projectID debe ser el ID de tu proyecto en GCP
	projectID := os.Getenv("ID_PROJECT")
	appEnv := os.Getenv("APP_ENV")

	// Crea el contexto
	ctx := context.Background()

	var client *secretmanager.Client
	var err error

	if appEnv == "production" {
		// En producción, usa la autenticación automática de App Engine
		client, err = secretmanager.NewClient(ctx)
	} else {
		// En desarrollo, usa un archivo de credenciales
		client, err = secretmanager.NewClient(ctx, option.WithCredentialsFile("credentials.json"))
	}

	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	// Construye el nombre del recurso del secreto
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretID)

	// Solicita el secreto al Secret Manager API
	result, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	})
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}

	// Devuelve el payload del secreto como string
	return string(result.Payload.Data), nil
}

func DBConnection() {
	// Obtén los secretos del Secret Manager
	host, err := getSecret("DB_HOST")
	if err != nil {
		log.Fatal(err)
	}

	user, err := getSecret("DB_USER")
	if err != nil {
		log.Fatal(err)
	}

	password, err := getSecret("DB_PASSWORD")
	if err != nil {
		log.Fatal(err)
	}

	dbName, err := getSecret("DB_NAME")
	if err != nil {
		log.Fatal(err)
	}

	port, err := getSecret("DB_PORT")
	if err != nil {
		log.Fatal(err)
	}

	// Construye la cadena de conexión y conecta a la base de datos
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, password, dbName, port)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Database connected")
	}
}
