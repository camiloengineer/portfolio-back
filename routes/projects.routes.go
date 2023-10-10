package routes

import (
	"encoding/json"
	"net/http"

	"github.com/camiloengineer/portfolio-back/db"
	"github.com/camiloengineer/portfolio-back/models"
)

func GetPracticePrjHandler(w http.ResponseWriter, r *http.Request) {
	var prj []models.Project
	db.DB.Find(&prj)
	json.NewEncoder(w).Encode(&prj)
}

func GetProfessionalPrjHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetProfessionalPrjHandler"))
}
