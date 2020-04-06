package controllers

import (
	"github.com/faeelol/lecture_recorder_backend/database"
	"github.com/faeelol/lecture_recorder_backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type UserController struct {
	ControllerBase
}

func (u *UserController) Test(c *gin.Context) {
	u.JsonSuccess(c, http.StatusOK, "test")
}

func (u *UserController) Create(c *gin.Context) {
	var newUser models.UserPost
	if err := c.Bind(&newUser); err != nil {
		u.JsonFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if !isUsernameFree(newUser.Username) {
		u.JsonFail(c, http.StatusBadRequest, "username is already taken")
		return
	}
	user := models.User{
		Username: newUser.Username,
		Password: Hash(newUser.Password),
	}
	if err := database.DB.Create(&user).Error; err != nil {
		panic(err)
	}
	u.JsonSuccess(c, http.StatusCreated, user.ToBasicUserSchema())
}

func (u *UserController) SelfInfo(c *gin.Context) {
	id := GetAuthUserClaims(c)
	var user models.User
	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		panic(err)
	}
	u.JsonSuccess(c, http.StatusOK, user.ToBasicUserSchema())
}

func isUsernameFree(username string) bool {
	var user models.User
	return gorm.IsRecordNotFoundError(database.DB.Where("username = ?", username).First(&user).Error)
}
