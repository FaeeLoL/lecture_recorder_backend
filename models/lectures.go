package models

import "github.com/jinzhu/gorm"

type Lecture struct {
	gorm.Model
	Name      string  `json:"name" form:"name"`
	AudioFile string  `json:"audio_file" form:"audio_file"`
	CourseId  uint    `json:"course_id" form:"course_id"`
	Notes     []Notes `json:"notes"`
}

type LecturePost struct {
	Name     string `json:"name" form:"name" binding:"required"`
	CourseId uint   `json:"course_id" form:"course_id" binding:"required"`
}

type LecturePut struct {
	Name *string `json:"name" form:"name"`
}

type Notes struct {
	gorm.Model
	LectureId uint   `json:"lecture_id"`
	Text      string `json:"text"`
	Picture   string `json:"picture"`
}

type BasicLectureSchema struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	AudioFile string `json:"audio_file"`
	CourseId  uint   `json:"course_id"`
}

func (l Lecture) ToBasicLectureSchema() BasicLectureSchema {
	return BasicLectureSchema{
		ID:        l.ID,
		Name:      l.Name,
		AudioFile: l.AudioFile,
		CourseId:  l.CourseId,
	}
}
