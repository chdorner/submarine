package data

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	return db, err
}

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&Settings{})
	if err != nil {
		panic(err)
	}
}
