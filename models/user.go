package models

import (
	"time"
)

type User struct {
	ID        uint   `gorm:"primary_key"`
	Username  string `gorm:"unique_index"`
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type UserPost struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type BasicUserSchema struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func (u User) ToBasicUserSchema() BasicUserSchema {
	return BasicUserSchema{
		ID:        u.ID,
		Username:  u.Username,
		CreatedAt: u.CreatedAt,
	}
}
