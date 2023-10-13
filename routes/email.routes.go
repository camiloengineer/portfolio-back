package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/option"
	"gopkg.in/mail.v2"

	"cloud.google.com/go/pubsub"
	"github.com/joho/godotenv"
)

var (
	projectID string
	topicID   string
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	projectID = os.Getenv("PROJECT_ID")
	topicID = os.Getenv("TOPIC_ID")
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

const developerEmailTemplate = `
<p>Name: {{.Name}},</p>
<p>Email: {{.Email}},</p>
<p>Message: {{.Message}}</p>
`

const userEmailTemplate = `
<p>Hello {{.Name}},</p>
<p>I hope you are well. I have received your email and will respond to it as soon as possible.</p>
<p>Best regards.</p>
`

func createBody(tmpl string, data interface{}) (string, error) {
	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return "", err
	}

	return body.String(), nil
}

func sendEmail(to string, subject string, body string) error {
	m := mail.NewMessage()
	m.SetHeader("From", "camilo@camiloengineer.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := mail.NewDialer("smtp.gmail.com", 587, "camilo@camiloengineer.com", "iowlgzkfcokjoqpr")

	return d.DialAndSend(m)
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

	// Wait for the result to get the message ID and error
	id, err := result.Get(ctx)
	if err != nil {
		log.Println("Error publishing message:", err)
		return err
	}

	log.Println("Published a message with ID:", id)
	return nil
}

func subscribeAndListenForMessages(ctx context.Context, subscriptionID string) error {
	client, err := createPubsubClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	sub := client.Subscription(subscriptionID)
	return sub.Receive(ctx, func(c context.Context, msg *pubsub.Message) { // Rename ctx to avoid shadowing
		var emailMessage EmailMessage
		if err := json.Unmarshal(msg.Data, &emailMessage); err != nil {
			log.Printf("Could not decode message data: %v", err)
			msg.Nack()
			return
		}

		bodyDeveloper, err := createBody(developerEmailTemplate, emailMessage)
		if err != nil {
			log.Printf("Error creating developer email body: %v", err)
			msg.Nack()
			return
		}

		if err := sendEmail("camilo@camiloengineer.com", emailMessage.Subject, bodyDeveloper); err != nil {
			log.Printf("Error sending email to developer: %v", err)
			msg.Nack()
			return
		}

		bodyUser, err := createBody(userEmailTemplate, emailMessage)
		if err != nil {
			log.Printf("Error creating user email body: %v", err)
			msg.Nack()
			return
		}

		if err := sendEmail(emailMessage.Email, "Thanks for contacting me!", bodyUser); err != nil {
			log.Printf("Error sending email to user: %v", err)
			msg.Nack()
			return
		}

		msg.Ack()
	})
}

func createPubsubClient(ctx context.Context) (*pubsub.Client, error) {
	appEnv := os.Getenv("APP_ENV")

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
