package models

import "github.com/jinzhu/gorm"

type Course struct {
	gorm.Model
	Name        string
	Description string
	Owner       uint
	Audios      []Lecture `gorm:"foreignkey:Audio"`
	Topic       uint
}

type BasicCourseSchema struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Audios      int    `json:"audios"`
	Topic       uint   `json:"topic"`
	IsOwner     bool   `json:"is_owner"`
}

type CoursePost struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type CoursePut struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (c Course) ToBasicCourseSchema(uid uint) *BasicCourseSchema {
	return &BasicCourseSchema{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		Audios:      len(c.Audios),
		Topic:       c.Topic,
		IsOwner:     c.Owner == uid,
	}
}
