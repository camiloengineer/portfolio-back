package main

import (
	"net/http"

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

	r.HandleFunc("/practiceprj", routes.GetPracticePrjHandler).Methods("GET")
	r.HandleFunc("/professionalprj", routes.GetProfessionalPrjHandler).Methods("GET")

	http.ListenAndServe(":3000", r)
}
