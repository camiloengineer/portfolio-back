package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/camiloengineer/portfolio-back/internal/email"
)

func SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	var emailData email.EmailMessage

	if err := json.NewDecoder(r.Body).Decode(&emailData); err != nil {
		log.Println("Error decoding the request body:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	log.Println("Request data decoded successfully:", emailData)

	if err := email.SendEmail(emailData); err != nil {
		log.Println("Error sending email:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println("Email sent successfully")

	// Responde al cliente
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email sent successfully"))
}
