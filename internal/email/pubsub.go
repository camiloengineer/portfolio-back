package email

import (
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/pubsub"
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
	projectID = config.Env.ProjectID
	topicID = config.Env.TopicID
	developerEmail = config.Env.DeveloperEmail
}

func SendEmail(emailData EmailMessage) error {
	msg, err := json.Marshal(emailData)
	if err != nil {
		return err
	}

	ctx := context.Background()
	if err := publishMessage(ctx, topicID, msg); err != nil {
		return err
	}

	return nil
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
		return err
	}
	defer client.Close()

	t := client.Topic(topicID)
	result := t.Publish(ctx, &pubsub.Message{
		Data: msg,
	})

	id, err := result.Get(ctx)
	if err != nil {
		return err
	}

	log.Println("Published a message with ID:", id)
	return nil
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

		developerBody := DeveloperBody
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

		userBody := UserBody
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
		log.Println("Email sent to user")

		msg.Ack()
	})
}
