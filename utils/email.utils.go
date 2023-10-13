package utils

import "gopkg.in/gomail.v2"

func sendEmail(to string, subject string, body string) error {
	d := gomail.NewDialer("smtp.gmail.com", 587, "your-email@gmail.com", "your-password")
	m := gomail.NewMessage()
	m.SetHeader("From", "your-email@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return d.DialAndSend(m)
}
