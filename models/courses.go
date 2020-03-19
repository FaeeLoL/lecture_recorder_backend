package models

import "github.com/jinzhu/gorm"

type Course struct {
	gorm.Model
	Name        string
	Description string
	Audios      []Audio `gorm:"foreignkey:Course"`
	Topic       uint
}

type CourseMinified struct {
	ID          uint
	Name        string
	Description string
	Audios      int
	Topic       uint
}

func (c Course) Minify() *CourseMinified {
	return &CourseMinified{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		Audios:      len(c.Audios),
		Topic:       c.Topic,
	}
}
