package models

import "gorm.io/gorm"

type ProjectCategories struct {
	gorm.Model

	ProjectID  uint `gorm:"not null"`
	CategoryID uint `gorm:"not null"`

	// Relaciones
	Project  Project
	Category Category
}
