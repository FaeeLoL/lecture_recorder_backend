package database

import (
	"fmt"
	"github.com/faeelol/lecture_recorder_backend/configs"
	"github.com/faeelol/lecture_recorder_backend/models"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(
		"mysql",
		fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			configs.DBConfigs.Login,
			configs.DBConfigs.Password,
			configs.DBConfigs.Address,
			configs.DBConfigs.Name,
		))
	if err != nil {
		return nil, err
	}
	db.DB().SetMaxIdleConns(100)
	DB = db
	db.AutoMigrate(&models.User{}, &models.Topic{}, &models.Course{}, &models.Lecture{}, &models.Note{})
	return db, err
}
