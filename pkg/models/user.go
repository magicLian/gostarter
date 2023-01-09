package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id        string         `json:"id" gorm:"column:id;primary_key;not null;type:varchar(255)"`
	Name      string         `json:"username"`
	IsAdmin   bool           `json:"isAdmin"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
}

type SignInUser struct {
	User
}
