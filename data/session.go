package data

import (
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	Token     string `gorm:"unique"`
	UserAgent string
	IP        string
}

type SessionCreate struct {
	UserAgent string
	IP        string
}
