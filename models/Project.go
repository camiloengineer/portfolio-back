package models

import "gorm.io/gorm"

type Project struct {
	gorm.Model

	Url            string `gorm:"size:255;not null"`
	Image          string `gorm:"size:255;not null"`
	IsProfessional bool   `gorm:"not null;default:false"`

	// Relaciones
	Translations []ProjectTranslation
	Categories   []Category `gorm:"many2many:project_categories;"`
}
