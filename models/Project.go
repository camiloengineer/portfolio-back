package models

import "gorm.io/gorm"

type Project struct {
	gorm.Model

	Url   string
	Image string
}
