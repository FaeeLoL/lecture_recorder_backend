package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	ControllerBase
}

func (u *UserController)Test(c *gin.Context) {
	u.JsonSuccess(c, http.StatusOK, "test")
}