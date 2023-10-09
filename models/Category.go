package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model

	Name string `gorm:"size:255;not null;unique"`

	// Relación
	Projects []Project `gorm:"many2many:project_categories;"`
}
