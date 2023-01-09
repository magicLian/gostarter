package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id         string         `json:"id"`
	Name       string         `json:"username"`
	CreateTime time.Time      `json:"createTime"`
	IsAdmin    bool           `json:"isAdmin"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"deletedAt"`
}

type SignInUser struct {
	User
}
