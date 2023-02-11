package data

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name string `gorm:"unique;not null;default:null;"`

	Bookmarks []Bookmark `gorm:"many2many:bookmark_tags;"`
}
