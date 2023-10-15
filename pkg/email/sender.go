package email

import (
	"bytes"
	"html/template"
	"log"
	"strconv"

	"gopkg.in/mail.v2"

	config "github.com/camiloengineer/portfolio-back/cmd"
)

var emailDialer *mail.Dialer
var DeveloperEmail string

func init() {
	host := config.Sec.AdminEmailHost
	port := config.Sec.AdminEmailPort

	intPort, err := strconv.Atoi(port)

	if err != nil {
		log.Fatalf("Error converting port to integer: %v", err)
	}

	emailDialer = mail.NewDialer(host, intPort, config.Sec.AdminEmailUser, config.Sec.AdminEmailPassword)
	DeveloperEmail = config.Env.DeveloperEmail

	log.Printf("Initialized emailDialer with host: %s, port: %d\n", host, intPort)
	log.Printf("DeveloperEmail set to: %s\n", DeveloperEmail)
}

func SendEmail(to, subject, body string) error {
	m := mail.NewMessage()
	m.SetHeader("From", DeveloperEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return emailDialer.DialAndSend(m)
}

func CreateBody(tmpl string, data interface{}) (string, error) {
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
