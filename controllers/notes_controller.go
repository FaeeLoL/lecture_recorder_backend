package controllers

import (
	"fmt"
	"github.com/faeelol/lecture_recorder_backend/database"
	"github.com/faeelol/lecture_recorder_backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
)

type NotesController struct {
	ControllerBase
}

func (n NotesController) Create(c *gin.Context) {
	var newNote models.NotePost

	err := c.Bind(&newNote)
	if err != nil {
		n.JsonFail(c, http.StatusBadRequest, err.Error())
		return
	}
	var lectureIdInt int
	if lectureId, res := c.Params.Get("lecture_id"); res {
		lectureIdInt, err = strconv.Atoi(lectureId)
		if err != nil {
			n.JsonFail(c, http.StatusBadRequest, fmt.Sprintf("Failed to convert lecture_id to int: %+v", err))
			return
		}
	} else {
		n.JsonFail(c, http.StatusBadRequest, "failed to get query paramL lecture_id")
		return
	}
	uid := GetAuthUserClaims(c)
	if !isLectureOwner(uid, uint(lectureIdInt)) {
		n.JsonFail(c, http.StatusForbidden, "the operation is forbidden")
		return
	}
	filename, err := SaveFile(c)
	if err != nil {
		n.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}

	if newNote.Text == "" && filename == "" {
		n.JsonFail(c, http.StatusBadRequest, "specify note `text` or `file`")
		return
	}

	note := models.Note{
		Text:      newNote.Text,
		Picture:   filename,
		LectureId: uint(lectureIdInt),
		Timestamp: newNote.Timestamp,
	}
	if err = database.DB.Save(&note).Error; err != nil {
		panic(err)
	}
	n.JsonSuccess(c, http.StatusCreated, note.ToBasicNoteSchema())
}

func (n NotesController) List(c *gin.Context) {
	var notes []models.Note
	var lectureIdInt int
	if lectureId, res := c.Params.Get("lecture_id"); res {
		var err error
		lectureIdInt, err = strconv.Atoi(lectureId)
		if err != nil {
			n.JsonFail(c, http.StatusBadRequest, fmt.Sprintf("Failed to convert lecture to int: %+v", err))
			return
		}
	} else {
		n.JsonFail(c, http.StatusBadRequest, "failed tp get query param lecture_id")
		return
	}
	if err := database.DB.Find(&notes, "lecture_id = ?", lectureIdInt).Error; err != nil {
		n.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	minifiedNotes := make([]models.BasicNoteSchema, 0)
	for _, note := range notes {
		minifiedNotes = append(minifiedNotes, note.ToBasicNoteSchema())
	}
	n.JsonSuccess(c, http.StatusOK, &minifiedNotes)
}

func (n NotesController) Get(c *gin.Context) {
	var note models.Note
	noteId, res := c.Params.Get("note_id")
	if !res {
		n.JsonFail(c, http.StatusBadRequest, "Empty note_id field")
		return
	}
	if err := database.DB.Find(&note, "id = ?", noteId).Error; err != nil {
		n.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	n.JsonSuccess(c, http.StatusOK, note.ToBasicNoteSchema())
}

func (n NotesController) Put(c *gin.Context) {
	var note models.Note

	id, res := c.Params.Get("note_id")
	if !res {
		n.JsonFail(c, http.StatusBadRequest, "Empty note_id field")
		return
	}
	if err := database.DB.Find(&note, "id = ?", id).Error; err != nil {
		n.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	uid := GetAuthUserClaims(c)
	if !isLectureOwner(uid, note.LectureId) {
		n.JsonFail(c, http.StatusForbidden, "the operation for current user is forbidden")
		return
	}

	filename, err := SaveFile(c)
	if err != nil {
		n.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	var changingNote models.NotePut
	err = c.Bind(&changingNote)
	if err != nil {
		n.JsonFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if changingNote.Text != nil {
		note.Text = *changingNote.Text
	}
	if changingNote.Timestamp != nil {
		note.Timestamp = *changingNote.Timestamp
	}
	if filename != "" {
		note.Picture = filename
	}

	if err := database.DB.Save(&note).Error; err != nil {
		panic(err)
	}
	n.JsonSuccess(c, http.StatusOK, note.ToBasicNoteSchema())
}

func (n NotesController) Delete(c *gin.Context) {
	var note models.Note

	id, res := c.Params.Get("note_id")
	if !res {
		n.JsonFail(c, http.StatusBadRequest, "Empty note_id field")
		return
	}
	if err := database.DB.Find(&note, "id = ?", id).Error; err != nil {
		n.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	uid := GetAuthUserClaims(c)
	if !isLectureOwner(uid, note.LectureId) {
		n.JsonFail(c, http.StatusForbidden, "the operation for current user is forbidden")
		return
	}

	if err := database.DB.Delete(&note, "id = ?", id).Error; err != nil {
		n.JsonFail(c, http.StatusInternalServerError, err.Error())
		return
	}
	n.JsonSuccess(c, http.StatusNoContent, nil)
}

func SaveFile(c *gin.Context) (string, error) {
	file, err := c.FormFile("file")
	filename := ""
	if err == nil {
		filename = newRandomFilename() + filepath.Ext(file.Filename)
		if err := c.SaveUploadedFile(file, path.Join("data", filename)); err != nil {
			return "", fmt.Errorf("faile save err: %s", err.Error())
		}
	}
	return filename, nil
}
