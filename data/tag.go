package data

import (
	"strings"

	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name        string `gorm:"unique;not null;default:null;"`
	DisplayName string

	Bookmarks []Bookmark `gorm:"many2many:bookmark_tags;"`
}

func (t *Tag) BeforeSave(tx *gorm.DB) error {
	t.Name = strings.ToLower(t.DisplayName)
	return nil
}
