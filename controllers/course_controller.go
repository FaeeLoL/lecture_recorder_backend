package controllers

import (
	"fmt"
	"github.com/faeelol/lecture_recorder_backend/database"
	"github.com/faeelol/lecture_recorder_backend/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type CourseController struct {
	ControllerBase
}

func (a *CourseController) Create(c *gin.Context) {
	var newCourse models.Course

	err := c.Bind(&newCourse)
	if err != nil {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if topicID, res := c.Params.Get("topic_id"); res {
		topicIDint, err := strconv.Atoi(topicID)
		if err != nil {
			log.Printf("Failed to convert topic_id to int: %+v", err)
			a.JsonFail(c, http.StatusBadRequest, fmt.Sprintf("Failed to convert topic_id to int: %+v", err))
			return
		}
		newCourse.Topic = uint(topicIDint)
	} else {
		log.Printf("Failed to get query param: topic_id")
		a.JsonFail(c, http.StatusBadRequest, "failed tp get query paramL topic_id")
		return
	}
	tx := database.DB.Begin()
	if err = tx.Error; err != nil {
		log.Printf("Failed to create transaction: %+v", err.Error())
		a.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err = tx.Save(&newCourse).Error; err != nil {
		tx.Rollback()
		a.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err := tx.Commit().Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	a.JsonSuccess(c, http.StatusCreated, newCourse.Minify())
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
	var minifiedCourses []models.CourseMinified
	for _, course := range courses {
		minifiedCourses = append(minifiedCourses, *course.Minify())
	}
	a.JsonSuccess(c, http.StatusOK, &minifiedCourses)
}

func (a *CourseController) Get(c *gin.Context) {
	var course models.Course
	courseID, res := c.Params.Get("course_id")
	if !res {
		a.JsonFail(c, http.StatusBadRequest, "Empty `course_id` path param")
		return
	}
	if err := database.DB.Find(&course, "id = ?", courseID).Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	a.JsonSuccess(c, http.StatusOK, course.Minify())
}

func (a *CourseController) Put(c *gin.Context) {
	var course models.Course

	courseID, res := c.Params.Get("course_id")
	if !res {
		a.JsonFail(c, http.StatusBadRequest, "Empty courseID field")
		return
	}
	if err := database.DB.Find(&course, "id = ?", courseID).Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	err := c.Bind(&course)
	if err != nil {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
		return
	}
	//todo
	id_i, _ := strconv.Atoi(courseID)
	if course.ID != uint(id_i) {
		a.JsonFail(c, http.StatusBadRequest, "Changing ID prohibited")
		return
	}
	if err := database.DB.Save(&course).Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	a.JsonSuccess(c, http.StatusOK, course.Minify())
}

func (a *CourseController) Delete(c *gin.Context) {
	var course models.Course

	courseID, res := c.Params.Get("course_id")
	if !res {
		a.JsonFail(c, http.StatusBadRequest, "Empty id field")
		return
	}
	if err := database.DB.Delete(&course, "id = ?", courseID).Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	a.JsonSuccess(c, http.StatusNoContent, nil)
}
