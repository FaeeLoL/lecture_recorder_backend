package models

import "github.com/jinzhu/gorm"

type Topic struct {
	gorm.Model
	Name        string
	Description string
	Courses     []Course `gorm:"foreignkey:Topic"`
}

type TopicMinified struct {
	ID          uint
	Name        string
	Description string
	Courses     int
}

func (t Topic) Minify() *TopicMinified {
	return &TopicMinified{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		Courses:     len(t.Courses),
	}
}
