package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/camiloengineer/portfolio-back/db"
	"github.com/camiloengineer/portfolio-back/models"
	"github.com/camiloengineer/portfolio-back/routes"
	"github.com/gorilla/mux"
)

func main() {

	db.DBConnection()

	db.DB.AutoMigrate(models.Category{})
	db.DB.AutoMigrate(models.Project{})
	db.DB.AutoMigrate(models.ProjectCategories{})
	db.DB.AutoMigrate(models.ProjectTranslation{})

	r := mux.NewRouter()

	r.HandleFunc("/", routes.HomeHandler)
	r.HandleFunc("/projects/innovation/{lang:[a-z]{2}}", routes.GetInnovationPrjHandler).Methods("GET")
	r.HandleFunc("/projects/professional/{lang:[a-z]{2}}", routes.GetProfessionalPrjHandler).Methods("GET")
	r.HandleFunc("/sendemail", routes.SendEmailHandler).Methods("POST")

	// Implementación del suscriptor en una goroutine.
	topicID := os.Getenv("TOPIC_ID")
	if topicID == "" {
		log.Fatal("Environment variable TOPIC_ID is not set")
	}

	go func() {
		ctx := context.Background()
		subscriptionID := topicID + "-sub"
		if err := routes.SubscribeAndListenForMessages(ctx, subscriptionID); err != nil {
			log.Printf("Error starting subscriber: %v", err)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(http.ListenAndServe(":"+port, r))
}
