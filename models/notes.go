package models

import (
	"github.com/jinzhu/gorm"
)

type Note struct {
	gorm.Model
	Text      string `json:"text"`
	Picture   string `json:"picture"`
	Timestamp uint   `json:"timestamp"`
	LectureId uint   `json:"lecture_id"`
}

type NotePost struct {
	Text      string `json:"text" form:"text"`
	Timestamp uint   `json:"timestamp" form:"timestamp" binding:"required"`
}

type NotePut struct {
	Text      *string `json:"text" form:"text"`
	Timestamp *uint   `json:"timestamp" form:"timestamp"`
}

type BasicNoteSchema struct {
	Id        uint   `json:"id"`
	Text      string `json:"text"`
	Picture   string `json:"picture"`
	LectureId uint   `json:"lecture_id"`
	Timestamp uint   `json:"timestamp"`
}

func (n Note) ToBasicNoteSchema() BasicNoteSchema {
	return BasicNoteSchema{
		Id:        n.ID,
		Text:      n.Text,
		Picture:   n.Picture,
		LectureId: n.LectureId,
		Timestamp: n.Timestamp,
	}
}
