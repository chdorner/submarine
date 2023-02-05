package data

import "gorm.io/gorm"

type Settings struct {
	gorm.Model
	Password string
}

type SettingsUpsert struct {
	Password string
}
