package controllers

import (
	"github.com/faeelol/lecture_recorder_backend/database"
	"github.com/faeelol/lecture_recorder_backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type TopicsController struct {
	ControllerBase
}

func (a *TopicsController) Create(c *gin.Context) {
	var newTopic models.TopicPost

	err := c.Bind(&newTopic)
	if err != nil {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if !isTopicNameUnique(newTopic.Name) {
		a.JsonFail(c, http.StatusBadRequest, "topic name is already taken")
		return
	}

	uid := GetAuthUserClaims(c)
	topic := models.Topic{
		Name:        newTopic.Name,
		Description: newTopic.Description,
		Owner:       uid,
		Courses:     nil,
	}
	if err := database.DB.Save(&topic).Error; err != nil {
		panic(err)
	}
	a.JsonSuccess(c, http.StatusCreated, topic.ToBasicTopicSchema(uid))
}

func (a *TopicsController) List(c *gin.Context) {
	var topics []models.Topic
	uid := c.Request.URL.Query().Get("user_id")
	if uid != "" {
		if err := database.DB.Where("owner = ?", uid).Find(&topics).Error; err != nil {
			a.JsonFail(c, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		if err := database.DB.Find(&topics).Error; err != nil {
			a.JsonFail(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	minifiedTopics := make([]models.BasicTopicSchema, 0)
	usID := GetAuthUserClaims(c)
	for _, topic := range topics {
		fixCoursesKey(&topic)
		minifiedTopics = append(minifiedTopics, *topic.ToBasicTopicSchema(usID))
	}
	a.JsonSuccess(c, http.StatusOK, &minifiedTopics)
}

func (a *TopicsController) Get(c *gin.Context) {
	uid := GetAuthUserClaims(c)
	var topic models.Topic
	topicID, res := c.Params.Get("topic_id")
	if !res {
		a.JsonFail(c, http.StatusBadRequest, "Empty topicID field")
		return
	}
	if err := database.DB.Find(&topic, "id = ?", topicID).Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	fixCoursesKey(&topic)
	a.JsonSuccess(c, http.StatusOK, topic.ToBasicTopicSchema(uid))
}

func fixCoursesKey(t *models.Topic) {
	database.DB.Where("topic = ?", t.ID).Find(&t.Courses)
}

func (a *TopicsController) Put(c *gin.Context) {
	var topic models.Topic

	id, res := c.Params.Get("topic_id")
	if !res {
		a.JsonFail(c, http.StatusBadRequest, "Empty id field")
		return
	}
	if err := database.DB.Find(&topic, "id = ?", id).Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	uid := GetAuthUserClaims(c)
	if topic.Owner != uid {
		a.JsonFail(c, http.StatusForbidden, "the operation for current user is forbidden")
		return
	}

	var changingTopic models.TopicPut
	err := c.Bind(&changingTopic)
	if err != nil {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if changingTopic.Name != nil {
		if !isTopicNameUnique(*changingTopic.Name) {
			a.JsonFail(c, http.StatusBadRequest, "the new topic name is already taken")
			return
		}
		topic.Name = *changingTopic.Name
	}
	if changingTopic.Description != nil {
		topic.Description = *changingTopic.Description
	}
	if err := database.DB.Save(&topic).Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	fixCoursesKey(&topic)
	a.JsonSuccess(c, http.StatusOK, topic.ToBasicTopicSchema(uid))
}

func (a *TopicsController) Delete(c *gin.Context) {
	var topic models.Topic

	id, res := c.Params.Get("topic_id")
	if !res {
		a.JsonFail(c, http.StatusBadRequest, "Empty id field")
		return
	}

	if err := database.DB.Find(&topic, "id = ?", id).Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}

	uid := GetAuthUserClaims(c)
	if topic.Owner != uid {
		a.JsonFail(c, http.StatusForbidden, "the operation for curent user is forbidden")
		return
	}

	if err := database.DB.Delete(&topic, "id = ?", id).Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	deleteTopicCourses(id)
	database.DB.Delete(&models.TopicSubscribe{}, "topic_id = ?", id)
	a.JsonSuccess(c, http.StatusNoContent, nil)
}

func deleteTopicCourses(topicID string) {
	var courses []models.Course
	database.DB.Where("topic = ?", topicID).Find(&courses)
	for _, course := range courses {
		database.DB.Delete(&models.CourseSubscribe{}, "course_id = ?", course.ID)
	}
	database.DB.Delete(models.Course{}, "topic = ?", topicID)
}

func isTopicNameUnique(name string) bool {
	var topic models.Topic
	return gorm.IsRecordNotFoundError(database.DB.Where("name = ?", name).First(&topic).Error)
}
