package controllers

import (
	"fmt"
	"github.com/faeelol/lecture_recorder_backend/database"
	"github.com/faeelol/lecture_recorder_backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"math/rand"
	"net/http"
	"path"
	"path/filepath"
)

type LecturesController struct {
	ControllerBase
}

func (l LecturesController) Create(c *gin.Context) {
	var newLecture models.LecturePost

	if err := c.Bind(&newLecture); err != nil {
		l.JsonFail(c, http.StatusBadRequest, err.Error())
		return
	}

	uid := GetAuthUserClaims(c)
	if !isCourseOwner(uid, newLecture.CourseId) {
		l.JsonFail(c, http.StatusForbidden, "the operation is forbidden")
		return
	}

	if !isLectureNameUnique(newLecture.Name, newLecture.CourseId) {
		l.JsonFail(c, http.StatusBadRequest, "lecture name already taken")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		l.JsonFail(c, http.StatusBadRequest, fmt.Sprintf("file form err: %s", err.Error()))
		return
	}
	filename := newRandomFilename() + filepath.Ext(file.Filename)

	if err := c.SaveUploadedFile(file, path.Join("data", filename)); err != nil {
		l.JsonFail(c, http.StatusBadRequest, fmt.Sprintf("faile same err: %s", err.Error()))
		return
	}

	lecture := models.Lecture{
		Name:      newLecture.Name,
		AudioFile: filename,
		CourseId:  newLecture.CourseId,
	}

	if err := database.DB.Save(&lecture).Error; err != nil {
		panic(err)
	}

	l.JsonSuccess(c, http.StatusCreated, lecture.ToBasicLectureSchema())
}

func (l LecturesController) Get(c *gin.Context) {
	var lecture models.Lecture
	lectureID, res := c.Params.Get("lecture_id")
	if !res {
		l.JsonFail(c, http.StatusBadRequest, "Empty lecture_id field")
		return
	}
	if err := database.DB.Find(&lecture, "id = ?", lectureID).Error; err != nil {
		l.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	l.JsonSuccess(c, http.StatusOK, lecture.ToBasicLectureSchema())
}

func (l LecturesController) Put(c *gin.Context) {
	var lecture models.Lecture

	id, res := c.Params.Get("lecture_id")
	if !res {
		l.JsonFail(c, http.StatusBadRequest, "Empty lecture_id field")
		return
	}
	if err := database.DB.Find(&lecture, "id = ?", id).Error; err != nil {
		l.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	uid := GetAuthUserClaims(c)
	if !isCourseOwner(uid, lecture.CourseId) {
		l.JsonFail(c, http.StatusForbidden, "the operation for current user is forbidden")
		return
	}

	var changingLecture models.LecturePut
	err := c.Bind(&changingLecture)
	if err != nil {
		l.JsonFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if changingLecture.Name != nil {
		if !isLectureNameUnique(*changingLecture.Name, lecture.CourseId) {
			l.JsonFail(c, http.StatusBadRequest, "the new lecture name is already taken")
			return
		}
		lecture.Name = *changingLecture.Name
	}
	if err := database.DB.Save(&lecture).Error; err != nil {
		l.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	l.JsonSuccess(c, http.StatusOK, lecture.ToBasicLectureSchema())
}

func (l LecturesController) Delete(c *gin.Context) {
	var lecture models.Lecture

	id, res := c.Params.Get("lecture_id")
	if !res {
		l.JsonFail(c, http.StatusBadRequest, "Empty lecture_id field")
		return
	}
	if err := database.DB.Find(&lecture, "id = ?", id).Error; err != nil {
		l.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	uid := GetAuthUserClaims(c)
	if !isCourseOwner(uid, lecture.CourseId) {
		l.JsonFail(c, http.StatusForbidden, "the operation for current user is forbidden")
		return
	}

	if err := database.DB.Delete(&lecture, "id = ?", id).Error; err != nil {
		l.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	l.JsonSuccess(c, http.StatusNoContent, nil)
}

func (l LecturesController) List(c *gin.Context) {
	courseId := c.Request.URL.Query().Get("course_id")
	var lectures []models.Lecture
	if err := database.DB.Where("course_id = ?", courseId).Find(&lectures).Error; err != nil {
		panic(err)
	}

	results := make([]models.BasicLectureSchema, 0)
	for _, lecture := range lectures {
		results = append(results, lecture.ToBasicLectureSchema())
	}
	l.JsonSuccess(c, http.StatusOK, results)
}

func isLectureNameUnique(name string, courseId uint) bool {
	var lecture models.Lecture
	return gorm.IsRecordNotFoundError(database.DB.Where("name = ? AND course_id = ?", name, courseId).First(&lecture).Error)
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const length = 16

func newRandomFilename() string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func isLectureOwner(uid uint, lectureId uint) bool {
	var lecture models.Lecture
	database.DB.Where("id = ?", lectureId).First(&lecture)
	return isCourseOwner(uid, lecture.CourseId)
}
