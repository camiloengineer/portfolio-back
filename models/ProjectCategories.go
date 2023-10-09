package models

import "gorm.io/gorm"

type ProjectCategories struct {
	gorm.Model

	Url   string
	Image string
}
