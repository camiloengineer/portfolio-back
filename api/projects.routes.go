package api

import (
	"encoding/json"
	"net/http"

	"github.com/camiloengineer/portfolio-back/internal/db"
	"github.com/camiloengineer/portfolio-back/internal/models"
	"github.com/gorilla/mux"
)

// @Summary Obtiene proyectos de innovación
// @Description Obtiene una lista de proyectos de innovación en función del idioma.
// @Tags Proyectos
// @Accept json
// @Produce json
// @Param lang path string true "Código de idioma (por ejemplo, 'es', 'en')"
// @Success 200 {array} ProjectResponse
// @Router /projects/innovation/{lang} [get]
func GetInnovationPrjHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var projects []models.Project
	var result []map[string]interface{}

	vars := mux.Vars(r)
	lang := vars["lang"]

	db.DB.Preload("Categories").Where("is_professional = ?", false).Find(&projects)
	for _, project := range projects {
		result = append(result, transformProjectToResponse(project, lang))
	}
	json.NewEncoder(w).Encode(result)
}

// @Summary Obtiene proyectos profesionales
// @Description Obtiene una lista de proyectos profesionales en función del idioma.
// @Tags Proyectos
// @Accept json
// @Produce json
// @Param lang path string true "Código de idioma (por ejemplo, 'es', 'en')"
// @Success 200 {array} ProjectResponse
// @Router /projects/professional/{lang} [get]
func GetProfessionalPrjHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var projects []models.Project
	var result []map[string]interface{}

	vars := mux.Vars(r)
	lang := vars["lang"]

	db.DB.Preload("Categories").Where("is_professional = ?", true).Find(&projects)
	for _, project := range projects {
		result = append(result, transformProjectToResponse(project, lang))
	}
	json.NewEncoder(w).Encode(result)
}

func transformProjectToResponse(project models.Project, lang string) map[string]interface{} {
	var categories []string
	for _, category := range project.Categories {
		categories = append(categories, category.Name)
	}

	translation := getTranslation(project.ID, lang)

	return map[string]interface{}{
		"title":       translation.Title,
		"category":    categories,
		"description": translation.Description,
		"url":         project.Url,
		"image":       project.Image,
	}
}

func getTranslation(projectID uint, lang string) models.ProjectTranslation {
	var translation models.ProjectTranslation
	db.DB.Where("project_id = ? AND language = ?", projectID, lang).First(&translation)
	return translation
}
