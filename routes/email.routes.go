package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"cloud.google.com/go/pubsub"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"gopkg.in/mail.v2"

	"github.com/camiloengineer/portfolio-back/cmd/config"
	"github.com/camiloengineer/portfolio-back/pkg/email"
)

var (
	projectID      string
	topicID        string
	developerEmail string
	emailDialer    *mail.Dialer
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	projectID = config.Env.ProjectID
	topicID = config.Env.TopicID
	developerEmail = config.Env.DeveloperEmail
}

func SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	var emailData EmailMessage

	if err := json.NewDecoder(r.Body).Decode(&emailData); err != nil {
		log.Println("Error decoding the request body:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	log.Println("Request data decoded successfully:", emailData)

	msg, err := json.Marshal(emailData)
	if err != nil {
		log.Println("Error marshaling email data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Publica el mensaje en Pub/Sub
	if err := publishMessage(context.Background(), topicID, msg); err != nil {
		log.Println("Error publishing message to Pub/Sub:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println("Message published successfully")

	// Responde al cliente
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email sent successfully"))
}

type EmailMessage struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
	Subject string `json:"subject"`
}

func publishMessage(ctx context.Context, topicID string, msg []byte) error {
	client, err := createPubsubClient(ctx)
	if err != nil {
		log.Println("Error creating Pubsub client:", err)
		return err
	}
	defer client.Close()

	t := client.Topic(topicID)
	result := t.Publish(ctx, &pubsub.Message{
		Data: msg,
	})

	id, err := result.Get(ctx)
	if err != nil {
		log.Println("Error publishing message:", err)
		return err
	}

	log.Println("Published a message with ID:", id)
	return nil
}

func SubscribeAndListenForMessages(ctx context.Context, subscriptionID string) error {

	client, err := createPubsubClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	sub := client.Subscription(subscriptionID)
	return sub.Receive(ctx, func(c context.Context, msg *pubsub.Message) {
		log.Println("Message received, ID:", msg.ID)

		var emailMessage EmailMessage
		if err := json.Unmarshal(msg.Data, &emailMessage); err != nil {
			log.Printf("Could not decode message data: %v", err)
			msg.Nack()
			return
		}

		log.Printf("Decoded message data: %+v", emailMessage)

		developerBody := email.DeveloperBody
		bodyDeveloper, err := email.CreateBody(developerBody, emailMessage)
		if err != nil {
			log.Printf("Error creating developer email body: %v", err)
			msg.Nack()
			return
		}

		if err := email.SendEmail("camilo@camiloengineer.com", emailMessage.Subject, bodyDeveloper); err != nil {
			log.Printf("Error sending email to developer: %v", err)
			msg.Nack()
			return
		}
		log.Println("Email sent to developer")

		userBody := email.DeveloperBody
		bodyUser, err := email.CreateBody(userBody, emailMessage)
		if err != nil {
			log.Printf("Error creating user email body: %v", err)
			msg.Nack()
			return
		}

		if err := email.SendEmail(emailMessage.Email, "Thanks for contacting me!", bodyUser); err != nil {
			log.Printf("Error sending email to user: %v", err)
			msg.Nack()
			return
		}
		log.Println("Email sent to user") // Log cuando el email al usuario es enviado

		msg.Ack()
	})
}

func createPubsubClient(ctx context.Context) (*pubsub.Client, error) {
	appEnv := config.Env.AppEnv

	var client *pubsub.Client
	var err error

	if appEnv == "production" {
		client, err = pubsub.NewClient(ctx, projectID)
		log.Println("Using production credentials for Pub/Sub.")
	} else {
		client, err = pubsub.NewClient(ctx, projectID, option.WithCredentialsFile("credentials.json"))
		log.Println("Using local credentials for Pub/Sub.")
	}

	if err != nil {
		return nil, err
	}
	return client, nil
}
