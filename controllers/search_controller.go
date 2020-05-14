package controllers

import (
	"github.com/faeelol/lecture_recorder_backend/database"
	"github.com/faeelol/lecture_recorder_backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SearchController struct {
	ControllerBase
}

func (s SearchController) Get(c *gin.Context) {
	uid := GetAuthUserClaims(c)
	key := c.Request.URL.Query().Get("key")
	if key == "" {
		s.JsonFail(c, http.StatusBadRequest, "Empty `key` field")
		return
	}
	var topics []models.Topic

	if err := database.DB.Where("name LIKE ?", "%"+key+"%").Find(&topics).Error; err != nil {
		panic(err)
	}
	minifiedTopics := make([]models.BasicTopicSchema, 0)
	for _, topic := range topics {
		minifiedTopics = append(minifiedTopics, *topic.ToBasicTopicSchema(uid))
	}
	var courses []models.Course
	if err := database.DB.Where("name LIKE ?", "%"+key+"%").Find(&courses).Error; err != nil {
		panic(err)
	}
	minifiedCourses := make([]models.BasicCourseSchema, 0)
	for _, course := range courses {
		minifiedCourses = append(minifiedCourses, *course.ToBasicCourseSchema(uid))
	}
	var lectures []models.Lecture
	if err := database.DB.Where("name LIKE ?", "%"+key+"%").Find(&lectures).Error; err != nil {
		panic(err)
	}
	minifiedLectures := make([]models.BasicLectureSchema, 0)
	for _, lecture := range lectures {
		minifiedLectures = append(minifiedLectures, lecture.ToBasicLectureSchema())
	}

	s.JsonSuccess(c, http.StatusOK, searchResponse{
		Topics:   minifiedTopics,
		Courses:  minifiedCourses,
		Lectures: minifiedLectures,
	})
}

type searchResponse struct {
	Topics   []models.BasicTopicSchema   `json:"Topics"`
	Courses  []models.BasicCourseSchema  `json:"courses"`
	Lectures []models.BasicLectureSchema `json:"lectures"`
}
