package config

import (
	"context"
	"fmt"
	"log"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type Secrets struct {
	DbHost             string
	DbName             string
	DbPassword         string
	DbPort             string
	DbUser             string
	AdminEmailHost     string
	AdminEmailPassword string
	AdminEmailPort     string
	AdminEmailUser     string
}

var Sec *Secrets

var (
	projectID string
	err       error
)

func init() {
	ctx := context.Background()
	client, err := createSecretClient(ctx)
	if err != nil {
		log.Fatalf("failed to create secretmanager client: %v", err)
	}

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	projectID = Env.ProjectID
	fmt.Printf("Project ID: %s\n", projectID)

	var (
		dbHost, dbName, dbPassword, dbPort, dbUser, adminEmailHost, adminEmailPassword, adminEmailPort, adminEmailUser string
	)

	// Manejar los errores devueltos por getSecret
	if dbHost, err = getSecret(ctx, client, "DB_HOST"); err != nil {
		log.Fatalf("Failed to retrieve DB_HOST: %v", err)
	}
	if dbName, err = getSecret(ctx, client, "DB_NAME"); err != nil {
		log.Fatalf("Failed to retrieve DB_NAME: %v", err)
	}
	if dbPassword, err = getSecret(ctx, client, "DB_PASSWORD"); err != nil {
		log.Fatalf("Failed to retrieve DB_PASSWORD: %v", err)
	}
	if dbPort, err = getSecret(ctx, client, "DB_PORT"); err != nil {
		log.Fatalf("Failed to retrieve DB_PORT: %v", err)
	}
	if dbUser, err = getSecret(ctx, client, "DB_USER"); err != nil {
		log.Fatalf("Failed to retrieve DB_USER: %v", err)
	}
	if adminEmailHost, err = getSecret(ctx, client, "ADMIN_EMAIL_HOST"); err != nil {
		log.Fatalf("Failed to retrieve ADMIN_EMAIL_HOST: %v", err)
	}
	if adminEmailPassword, err = getSecret(ctx, client, "ADMIN_EMAIL_PASSWORD"); err != nil {
		log.Fatalf("Failed to retrieve ADMIN_EMAIL_PASSWORD: %v", err)
	}
	if adminEmailPort, err = getSecret(ctx, client, "ADMIN_EMAIL_PORT"); err != nil {
		log.Fatalf("Failed to retrieve ADMIN_EMAIL_PORT: %v", err)
	}
	if adminEmailUser, err = getSecret(ctx, client, "ADMIN_EMAIL_USER"); err != nil {
		log.Fatalf("Failed to retrieve ADMIN_EMAIL_USER: %v", err)
	}

	Sec = &Secrets{
		DbHost:             dbHost,
		DbName:             dbName,
		DbPassword:         dbPassword,
		DbPort:             dbPort,
		DbUser:             dbUser,
		AdminEmailHost:     adminEmailHost,
		AdminEmailPassword: adminEmailPassword,
		AdminEmailPort:     adminEmailPort,
		AdminEmailUser:     adminEmailUser,
	}
}

func createSecretClient(ctx context.Context) (*secretmanager.Client, error) {
	appEnv := Env.AppEnv

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
		return nil, err
	}
	return client, nil
}

func getSecret(ctx context.Context, client *secretmanager.Client, secretID string) (string, error) {
	log.Printf("Getting secret with ID: %s", secretID)
	// Build the request.
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretID)

	log.Printf("Constructed secret name: %s\n", name)
	// Call the API.
	result, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	})
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}

	return string(result.Payload.Data), nil
}
