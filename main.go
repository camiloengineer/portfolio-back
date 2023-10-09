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

	http.ListenAndServe(":3000", r)
}
