// cmd/config/secret.go

package config

import (
	"context"
	"log"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"google.golang.org/api/option"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type Secrets struct {
	DbUser     string
	DbPassword string
	SmtpUser   string
	SmtpPass   string
}

var Sec *Secrets

func init() {
	ctx := context.Background()
	client, err := createSecretClient(ctx)
	if err != nil {
		log.Fatalf("failed to create secretmanager client: %v", err)
	}

	Sec = &Secrets{
		DbUser:     getSecret(ctx, client, "DB_USER"),
		DbPassword: getSecret(ctx, client, "DB_PASSWORD"),
		SmtpUser:   getSecret(ctx, client, "SMTP_USER"),
		SmtpPass:   getSecret(ctx, client, "SMTP_PASS"),
	}
}

func createSecretClient(ctx context.Context) (*secretmanager.Client, error) {
	creds, exists := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if !exists {
		log.Println("No GOOGLE_APPLICATION_CREDENTIALS provided, using default")
		return secretmanager.NewClient(ctx)
	}
	return secretmanager.NewClient(ctx, option.WithCredentialsFile(creds))
}

func getSecret(ctx context.Context, client *secretmanager.Client, secretID string) string {
	// Build the request.
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/" + Env.ProjectID + "/secrets/" + secretID + "/versions/latest",
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		log.Printf("failed to access secret version: %v", err)
		return ""
	}

	// Return the secret payload as a string.
	return string(result.Payload.Data)
}
