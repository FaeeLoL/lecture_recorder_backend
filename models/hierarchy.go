package models

import "github.com/jinzhu/gorm"

type Topic struct {
	gorm.Model
	Name string
	Description string
	Courses []Course `gorm:"foreignkey:Topic"`
}

type Course struct {
	gorm.Model
	Name string
	Description string
	Audios []Audio `gorm:"foreignkey:Course"`
	Topic uint
}

type Audio struct {
	gorm.Model
	Name string
	Course uint
}