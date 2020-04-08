package controllers

import (
	"fmt"
	"github.com/faeelol/lecture_recorder_backend/database"
	"github.com/faeelol/lecture_recorder_backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
)

type CourseController struct {
	ControllerBase
}

func (a *CourseController) Create(c *gin.Context) {
	var newCourse models.CoursePost

	err := c.Bind(&newCourse)
	if err != nil {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
		return
	}
	var topicIDint int
	if topicID, res := c.Params.Get("topic_id"); res {
		topicIDint, err = strconv.Atoi(topicID)
		if err != nil {
			log.Printf("Failed to convert topic_id to int: %+v", err)
			a.JsonFail(c, http.StatusBadRequest, fmt.Sprintf("Failed to convert topic_id to int: %+v", err))
			return
		}
	} else {
		log.Printf("Failed to get query param: topic_id")
		a.JsonFail(c, http.StatusBadRequest, "failed tp get query paramL topic_id")
		return
	}
	uid := GetAuthUserClaims(c)
	if !isTopicOwner(uid, uint(topicIDint)) {
		a.JsonFail(c, http.StatusForbidden, "the operation is forbidden")
		return
	}
	course := models.Course{
		Name:        newCourse.Name,
		Description: newCourse.Description,
		Owner:       uid,
		Topic:       uint(topicIDint),
	}
	if err = database.DB.Save(&course).Error; err != nil {
		panic(err)
	}
	a.JsonSuccess(c, http.StatusCreated, course.ToBasicCourseSchema())
}

func (a *CourseController) List(c *gin.Context) {
	var courses []models.Course
	var topicIDint int
	if topicID, res := c.Params.Get("topic_id"); res {
		var err error
		topicIDint, err = strconv.Atoi(topicID)
		if err != nil {
			log.Printf("Failed to convert topic_id to int: %+v", err)
			a.JsonFail(c, http.StatusBadRequest, fmt.Sprintf("Failed to convert topic_id to int: %+v", err))
			return
		}
	} else {
		log.Printf("Failed to get query param: topic_id")
		a.JsonFail(c, http.StatusBadRequest, "failed tp get query paramL topic_id")
		return
	}
	if err := database.DB.Find(&courses, "Topic = ?", topicIDint).Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	var minifiedCourses []models.BasicCourseSchema
	for _, course := range courses {
		minifiedCourses = append(minifiedCourses, *course.ToBasicCourseSchema())
	}
	a.JsonSuccess(c, http.StatusOK, &minifiedCourses)
}

func (a *CourseController) Get(c *gin.Context) {
	var course models.Course
	var courseIDint int
	if courseID, res := c.Params.Get("course_id"); res {
		var err error
		courseIDint, err = strconv.Atoi(courseID)
		if err != nil {
			log.Printf("Failed to convert course to int: %+v", err)
			a.JsonFail(c, http.StatusBadRequest, fmt.Sprintf("Failed to convert course to int: %+v", err))
			return
		}
	} else {
		log.Printf("Failed to get query param: topic_id")
		a.JsonFail(c, http.StatusBadRequest, "failed tp get query paramL course_id")
		return
	}
	if err := database.DB.Find(&course, "id = ?", courseIDint).Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	a.JsonSuccess(c, http.StatusOK, course.ToBasicCourseSchema())
}

func (a *CourseController) Put(c *gin.Context) {
	var course models.Course

	var courseIDint int
	if courseID, res := c.Params.Get("course_id"); res {
		var err error
		courseIDint, err = strconv.Atoi(courseID)
		if err != nil {
			log.Printf("Failed to convert course to int: %+v", err)
			a.JsonFail(c, http.StatusBadRequest, fmt.Sprintf("Failed to convert course to int: %+v", err))
			return
		}
	} else {
		log.Printf("Failed to get query param: topic_id")
		a.JsonFail(c, http.StatusBadRequest, "failed tp get query paramL course_id")
		return
	}

	uid := GetAuthUserClaims(c)

	if err := database.DB.Find(&course, "id = ?", courseIDint).Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	if course.Owner != uid {
		a.JsonFail(c, http.StatusForbidden, "the operation is forbidden")
		return
	}
	var modifiedCourse models.CoursePut
	err := c.Bind(&modifiedCourse)
	if err != nil {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if modifiedCourse.Name != nil {
		if !isCourseNameUnique(course.Topic, *modifiedCourse.Name) {
			a.JsonFail(c, http.StatusBadRequest, "the course name is already taken")
			return
		}
		course.Name = *modifiedCourse.Name
	}
	if modifiedCourse.Description != nil {
		course.Description = *modifiedCourse.Description
	}
	if err := database.DB.Save(&course).Error; err != nil {
		panic(err)
	}
	a.JsonSuccess(c, http.StatusOK, course.ToBasicCourseSchema())
}

func (a *CourseController) Delete(c *gin.Context) {
	var course models.Course

	var courseIDint int
	if courseID, res := c.Params.Get("course_id"); res {
		var err error
		courseIDint, err = strconv.Atoi(courseID)
		if err != nil {
			log.Printf("Failed to convert course to int: %+v", err)
			a.JsonFail(c, http.StatusBadRequest, fmt.Sprintf("Failed to convert course to int: %+v", err))
			return
		}
	} else {
		log.Printf("Failed to get query param: topic_id")
		a.JsonFail(c, http.StatusBadRequest, "failed tp get query paramL course_id")
		return
	}

	uid := GetAuthUserClaims(c)
	if !isCourseOwner(uid, uint(courseIDint)) {
		a.JsonFail(c, http.StatusForbidden, "the operation is forbidden")
		return
	}
	if err := database.DB.Delete(&course, "id = ?", courseIDint).Error; err != nil {
		panic(err)
	}
	a.JsonSuccess(c, http.StatusNoContent, nil)
}

func isTopicOwner(userId uint, topicId uint) bool {
	var topic models.Topic
	fmt.Println(userId, topicId)
	if err := database.DB.Where("id = ?", topicId).Find(&topic).Error; err != nil {
		panic(err)
	}
	return topic.Owner == userId
}

func isCourseOwner(userId uint, courseId uint) bool {
	var course models.Course
	if err := database.DB.Where("id = ?", courseId).First(&course).Error; err != nil {
		return false
	}
	return course.Owner == userId
}

func isCourseNameUnique(courseId uint, name string) bool {
	var course models.Course
	return gorm.IsRecordNotFoundError(
		database.DB.Where("id = ? AND name = ?", courseId, name).Find(&course).Error)
}
