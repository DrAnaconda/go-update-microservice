package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name          string `gorm:"uniqueIndex"`
	PathToProduct string `gorm:"uniqueIndex"`
	Password      string
}
