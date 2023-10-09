package models

import "gorm.io/gorm"

type ProjectTranslation struct {
	gorm.Model

	ProjectID   uint   `gorm:"not null;index"`
	Language    string `gorm:"size:2;not null;index"`
	Title       string `gorm:"size:255;not null"`
	Description string `gorm:"type:text;not null"`
	ButtonText  string `gorm:"size:255;not null"`

	// Relación
	Project Project
}
