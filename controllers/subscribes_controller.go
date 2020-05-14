package controllers

import (
	"github.com/faeelol/lecture_recorder_backend/database"
	"github.com/faeelol/lecture_recorder_backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

type SubscribesController struct {
	ControllerBase
}

func isCourseExists(id uint) bool {
	var course models.Course
	err := database.DB.Where("id = ?", id).Find(&course).Error
	return !gorm.IsRecordNotFoundError(err)
}

func isTopicExists(id uint) bool {
	var topic models.Topic
	err := database.DB.Where("id = ?", id).Find(&topic).Error
	return !gorm.IsRecordNotFoundError(err)
}

func (s SubscribesController) Create(c *gin.Context) {
	uid := GetAuthUserClaims(c)
	if topicID := c.Request.URL.Query().Get("topic_id"); topicID != "" {
		topicIDint, err := strconv.Atoi(topicID)
		if err != nil {
			s.JsonFail(c, http.StatusBadRequest, "Invalid topic_id format")
			return
		}
		if !isTopicExists(uint(topicIDint)) {
			s.JsonFail(c, http.StatusBadRequest, "Topic with such id not exists")
			return
		}
		topicSubscribe := models.TopicSubscribe{
			UserID:  uid,
			TopicID: uint(topicIDint),
		}
		if err := database.DB.Where("user_id = ? AND topic_id = ?",
			uid, uint(topicIDint)).Find(&topicSubscribe).Error;
			!gorm.IsRecordNotFoundError(err) {
			s.JsonFail(c, http.StatusBadRequest, "Subscription for this topic already done")
			return
		}
		if err := database.DB.Create(&topicSubscribe).Error; err != nil {
			panic(err)
		}
		s.JsonSuccess(c, http.StatusCreated, topicSubscribe)
		return
	} else if courseID := c.Request.URL.Query().Get("course_id"); courseID != "" {
		courseIDint, err := strconv.Atoi(courseID)
		if err != nil {
			s.JsonFail(c, http.StatusBadRequest, "Invalid topic_id format")
			return
		}
		if !isCourseExists(uint(courseIDint)) {
			s.JsonFail(c, http.StatusBadRequest, "Course with such id not exists")
			return
		}
		courseSubscribe := models.CourseSubscribe{
			UserID:   uid,
			CourseID: uint(courseIDint),
		}
		if err := database.DB.Where("user_id = ? AND course_id = ?",
			uid, uint(courseIDint)).Find(&courseSubscribe).Error;
			!gorm.IsRecordNotFoundError(err) {
			s.JsonFail(c, http.StatusBadRequest, "Subscription for this course already done")
			return
		}
		if err := database.DB.Create(&courseSubscribe).Error; err != nil {
			panic(err)
		}
		s.JsonSuccess(c, http.StatusCreated, courseSubscribe)
		return
	}
	s.JsonFail(c, http.StatusBadRequest, "Please specify `topic_id` or `course_id in the params")
}

func (s SubscribesController) Delete(c *gin.Context) {
	uid := GetAuthUserClaims(c)
	if topicID := c.Request.URL.Query().Get("topic_id"); topicID != "" {
		topicIDint, err := strconv.Atoi(topicID)
		if err != nil {
			s.JsonFail(c, http.StatusBadRequest, "Invalid topic_id format")
			return
		}
		if !isTopicExists(uint(topicIDint)) {
			s.JsonFail(c, http.StatusBadRequest, "Topic with such id not exists")
			return
		}
		topicSubscribe := models.TopicSubscribe{
			UserID:  uid,
			TopicID: uint(topicIDint),
		}
		if err := database.DB.Delete(&topicSubscribe, "user_id = ? AND topic_id = ?", uid, uint(topicIDint)).Error; err != nil {
			panic(err)
		}
		s.JsonSuccess(c, http.StatusNoContent, topicSubscribe)
		return
	} else if courseID := c.Request.URL.Query().Get("course_id"); courseID != "" {
		courseIDint, err := strconv.Atoi(courseID)
		if err != nil {
			s.JsonFail(c, http.StatusBadRequest, "Invalid topic_id format")
			return
		}
		if !isCourseExists(uint(courseIDint)) {
			s.JsonFail(c, http.StatusBadRequest, "Course with such id not exists")
			return
		}
		courseSubscribe := models.CourseSubscribe{
			UserID:   uid,
			CourseID: uint(courseIDint),
		}
		if err := database.DB.Delete(&courseSubscribe, "user_id = ? AND course_id = ?", uid, uint(courseIDint)).Error; err != nil {
			panic(err)
		}
		s.JsonSuccess(c, http.StatusNoContent, courseSubscribe)
		return
	}
	s.JsonFail(c, http.StatusBadRequest, "Please specify `topic_id` or `course_id in the params")
}

func (s SubscribesController) List(c *gin.Context) {
	uid := GetAuthUserClaims(c)
	var topicSubs []models.TopicSubscribe
	if err := database.DB.Where("user_id = ?", uid).Find(&topicSubs).Error; err != nil {
		panic(err)
	}

	topics := make([]models.Topic, 0)
	for _, s := range topicSubs {
		var t models.Topic
		if err := database.DB.Where("id = ?", s.TopicID).First(&t).Error; err != nil {
			panic(err)
		}
		topics = append(topics, t)
	}
	topicsBasic := make([]models.BasicTopicSchema, 0)
	for _, t := range topics {
		fixCoursesKey(&t)
		topicsBasic = append(topicsBasic, *t.ToBasicTopicSchema())
	}
	var courseSubs []models.CourseSubscribe
	if err := database.DB.Where("user_id = ?", uid).Find(&courseSubs).Error; err != nil {
		panic(err)
	}
	courses := make([]models.Course, 0)
	for _, s := range courseSubs {
		var c models.Course
		if err := database.DB.Where("id = ?", s.CourseID).First(&c).Error; err != nil {
			panic(err)
		}
		courses = append(courses, c)
	}
	coursesBasic := make([]models.BasicCourseSchema, 0)
	for _, c := range courses {
		fixLecturesKey(&c)
		coursesBasic = append(coursesBasic, *c.ToBasicCourseSchema())
	}
	s.JsonSuccess(c, http.StatusOK, subsReturn{
		Topics:  topicsBasic,
		Courses: coursesBasic,
	})
}

type subsReturn struct {
	Topics  []models.BasicTopicSchema  `json:"topics"`
	Courses []models.BasicCourseSchema `json:"courses"`
}
