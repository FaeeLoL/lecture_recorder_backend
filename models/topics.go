package models

import (
	"github.com/jinzhu/gorm"
)

type Topic struct {
	gorm.Model
	Name        string
	Description string
	Owner       uint
	Courses     []Course `gorm:"foreignkey:Courses"`
}

type TopicPost struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type TopicPut struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type BasicTopicSchema struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Courses     int    `json:"courses"`
}

func (t Topic) ToBasicTopicSchema() *BasicTopicSchema {
	return &BasicTopicSchema{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		Courses:     len(t.Courses), //todo fix relations
	}
}
