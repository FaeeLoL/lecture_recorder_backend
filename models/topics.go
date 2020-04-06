package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Topic struct {
	gorm.Model
	Name        string
	Description string
	Courses     []Course `gorm:"foreignkey:TopicRefer"`
}

type TopicMinified struct {
	ID          uint
	Name        string
	Description string
	Courses     int
}

func (t Topic) Minify() *TopicMinified {
	for _, c := range t.Courses {
		fmt.Printf("courses: %+v\n", c)
	}
	return &TopicMinified{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		Courses:     len(t.Courses),
	}
}
