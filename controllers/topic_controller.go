package controllers

import (
	"github.com/faeelol/lecture_recorder_backend/database"
	"github.com/faeelol/lecture_recorder_backend/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type TopicsController struct {
	ControllerBase
}

func (a *TopicsController) Create(c *gin.Context) {
	var newTopic models.Topic

	err := c.Bind(&newTopic)
	if err != nil {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
		return
	}
	tx := database.DB.Begin()
	if err = tx.Error; err != nil {
		log.Printf("Failed to create transaction: %+v", err.Error())
		a.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err = tx.Save(&newTopic).Error; err != nil {
		tx.Rollback()
		a.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err := tx.Commit().Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	a.JsonSuccess(c, http.StatusCreated, newTopic.Minify())
}

func (a *TopicsController) List(c *gin.Context) {
	var topics []models.Topic
	if err := database.DB.Find(&topics).Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	var minifiedTopics []models.TopicMinified
	for _, topic := range topics {
		minifiedTopics = append(minifiedTopics, *topic.Minify())
	}
	a.JsonSuccess(c, http.StatusOK, &minifiedTopics)
}

func (a *TopicsController) Get(c *gin.Context) {
	var topic models.Topic
	id, res := c.Params.Get("id")
	if !res {
		a.JsonFail(c, http.StatusBadRequest, "Empty id field")
		return
	}
	if err := database.DB.Find(&topic, "id = ?", id).Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	a.JsonSuccess(c, http.StatusOK, topic.Minify())
}

func (a *TopicsController) Put(c *gin.Context) {
	var topic models.Topic

	id, res := c.Params.Get("id")
	if !res {
		a.JsonFail(c, http.StatusBadRequest, "Empty id field")
		return
	}
	if err := database.DB.Find(&topic, "id = ?", id).Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	err := c.Bind(&topic)
	if err != nil {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
		return
	}
	//todo
	id_i, _ := strconv.Atoi(id)
	if topic.ID != uint(id_i) {
		a.JsonFail(c, http.StatusBadRequest, "Changing ID prohibited")
		return
	}
	if err := database.DB.Save(&topic).Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	a.JsonSuccess(c, http.StatusOK, topic.Minify())
}

func (a *TopicsController) Delete(c *gin.Context) {
	var topic models.Topic

	id, res := c.Params.Get("id")
	if !res {
		a.JsonFail(c, http.StatusBadRequest, "Empty id field")
		return
	}
	if err := database.DB.Delete(&topic, "id = ?", id).Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	a.JsonSuccess(c, http.StatusNoContent, nil)
}
