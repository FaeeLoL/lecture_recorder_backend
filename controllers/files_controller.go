package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

type FilesController struct {
	ControllerBase
}

func (l FilesController) GetFile(c *gin.Context) {
	filename, res := c.Params.Get("file")
	if !res {
		l.JsonFail(c, http.StatusBadRequest, "Empty file field")
		return
	}
	c.File(path.Join("data", filename))
}

