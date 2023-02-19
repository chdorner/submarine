package data

import (
	"strings"

	"github.com/glebarez/sqlite"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Connect(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	return db, err
}

func Migrate(db *gorm.DB) {
	m := NewMigrator(db)
	err := m.Migrate()
	if err != nil {
		panic(err)
	}
}

func NewMigrator(db *gorm.DB) *gormigrate.Gormigrate {
	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		// initial migration
		{
			ID: "202302191130",
			Migrate: func(tx *gorm.DB) error {
				type Settings struct {
					gorm.Model
					Password string
				}
				type Session struct {
					gorm.Model
					Token     string `gorm:"unique"`
					UserAgent string
					IP        string
				}
				type Bookmark struct {
					gorm.Model
					URL         string `gorm:"not null;default:null"`
					Title       string
					Description string
					Privacy     BookmarkPrivacy `gorm:"default:'private'"`

					Tags []Tag `gorm:"many2many:bookmark_tags;"`
				}
				type Tag struct {
					gorm.Model
					Name string `gorm:"unique;not null;default:null;"`

					Bookmarks []Bookmark `gorm:"many2many:bookmark_tags;"`
				}
				return tx.AutoMigrate(
					&Settings{},
					&Session{},
					&Bookmark{},
					&Tag{},
				)
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(
					"settings",
					"sessions",
					"bookmarks",
					"tags",
					"bookmark_tags",
				)
			},
		},
		{
			ID: "202302191830",
			Migrate: func(tx *gorm.DB) error {
				type Tag struct {
					gorm.Model
					Name        string `gorm:"unique;not null;default:null;"`
					DisplayName string

					Bookmarks []Bookmark `gorm:"many2many:bookmark_tags;"`
				}
				err := tx.AutoMigrate(&Tag{})
				if err != nil {
					return err
				}

				var tags []Tag
				tx.Find(&tags)
				for _, tag := range tags {
					tag.DisplayName = tag.Name
					tag.Name = strings.ToLower(tag.Name)
					tx.Save(tag)
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				type Tag struct {
					gorm.Model
					Name        string `gorm:"unique;not null;default:null;"`
					DisplayName string

					Bookmarks []Bookmark `gorm:"many2many:bookmark_tags;"`
				}

				var tags []Tag
				tx.Find(&tags)
				for _, tag := range tags {
					tag.Name = tag.DisplayName
					tx.Save(tag)
				}

				return tx.Migrator().DropColumn(&Tag{}, "DisplayName")
			},
		},
	})
}
