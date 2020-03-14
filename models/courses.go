package models

import "github.com/jinzhu/gorm"

type Course struct {
	gorm.Model
	Name        string
	Description string
	Audios      []Audio `gorm:"foreignkey:Course"`
	Topic       uint
}