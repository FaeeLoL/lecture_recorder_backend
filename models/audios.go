package models

import "github.com/jinzhu/gorm"

type Audio struct {
	gorm.Model
	Name   string
	Course uint
}
