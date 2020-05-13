package models

import "github.com/jinzhu/gorm"

type TopicSubscribe struct {
	gorm.Model
	UserID  uint
	TopicID uint
}

type BasicTopicSubscribeSchema struct {
	ID      uint `json:"id"`
	UserID  uint `json:"user_id"`
	TopicID uint `json:"topic_id"`
}

func (t TopicSubscribe) ToBasicTopicSubscribeSchema() *BasicTopicSubscribeSchema {
	return &BasicTopicSubscribeSchema{
		ID:      t.ID,
		UserID:  t.UserID,
		TopicID: t.TopicID,
	}
}

type CourseSubscribe struct {
	gorm.Model
	UserID   uint
	CourseID uint
}

type BasicCourseSubscribeSchema struct {
	ID       uint `json:"id"`
	UserID   uint `json:"user_id"`
	CourseID uint `json:"course_id"`
}

func (t CourseSubscribe) ToBasicCourseSubscribeSchema() *BasicCourseSubscribeSchema {
	return &BasicCourseSubscribeSchema{
		ID:      t.ID,
		UserID:  t.UserID,
		CourseID: t.CourseID,
	}
}
