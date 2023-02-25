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
	return gormigrate.New(db, &gormigrate.Options{
		TableName:      "migrations",
		IDColumnName:   "id",
		IDColumnSize:   255,
		UseTransaction: true,
	}, []*gormigrate.Migration{
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
		{
			ID: "202302251200",
			Migrate: func(tx *gorm.DB) error {
				err := tx.Exec(`CREATE VIRTUAL TABLE tags_fts USING fts5(
					display_name,
					content="tags",
					content_rowid="id"
				);`).Error
				if err != nil {
					return err
				}

				err = tx.Exec(`CREATE TRIGGER tags_ai AFTER INSERT ON tags BEGIN
					INSERT INTO tags_fts(rowid, display_name) VALUES (new.id, new.display_name);
				END;`).Error
				if err != nil {
					return err
				}

				err = tx.Exec(`CREATE TRIGGER tags_ad AFTER DELETE ON tags BEGIN
					INSERT INTO tags_fts(tags_fts, rowid, display_name) VALUES('delete', old.id, old.display_name);
				END;`).Error
				if err != nil {
					return err
				}

				err = tx.Exec(`CREATE TRIGGER tags_au AFTER UPDATE ON tags BEGIN
					INSERT INTO tags_fts(tags_fts, rowid, display_name) VALUES('delete', old.id, old.display_name);
					INSERT INTO tags_fts(rowid, display_name) VALUES (new.id, new.display_name);
				END;`).Error
				if err != nil {
					return err
				}

				err = tx.Exec(`INSERT INTO tags_fts(rowid, display_name) SELECT id, display_name from tags;`).Error
				if err != nil {
					return err
				}

				err = tx.Exec(`CREATE VIRTUAL TABLE bookmarks_fts USING fts5(
					url,
					title,
					description,
					content="bookmarks",
					content_rowid="id"
				);`).Error
				if err != nil {
					return err
				}

				err = tx.Exec(`CREATE TRIGGER bookmarks_ai AFTER INSERT ON bookmarks BEGIN
					INSERT INTO bookmarks_fts(rowid, url, title, description) VALUES (new.id, new.url, new.title, new.description);
				END;`).Error
				if err != nil {
					return err
				}

				err = tx.Exec(`CREATE TRIGGER bookmarks_ad AFTER DELETE ON bookmarks BEGIN
					INSERT INTO bookmarks_fts(bookmarks_fts, rowid, url, title, description) VALUES('delete', old.id, old.url, old.title, old.description);
				END;`).Error
				if err != nil {
					return err
				}

				err = tx.Exec(`CREATE TRIGGER bookmarks_au AFTER UPDATE ON bookmarks BEGIN
					INSERT INTO bookmarks_fts(bookmarks_fts, rowid, url, title, description) VALUES('delete', old.id, old.url, old.title, old.description);
					INSERT INTO bookmarks_fts(rowid, url, title, description) VALUES (new.id, new.url, new.title, new.description);
				END;`).Error
				if err != nil {
					return err
				}

				return tx.Exec(`INSERT INTO bookmarks_fts(rowid, url, title, description) SELECT id, url, title, description from bookmarks;`).Error
			},
			Rollback: func(tx *gorm.DB) error {
				err := tx.Exec("DROP TABLE tags_fts;").Error
				if err != nil {
					return err
				}

				err = tx.Exec("DROP TRIGGER tags_ai;").Error
				if err != nil {
					return err
				}

				err = tx.Exec("DROP TRIGGER tags_ad;").Error
				if err != nil {
					return err
				}

				err = tx.Exec("DROP TRIGGER tags_au;").Error
				if err != nil {
					return err
				}

				err = tx.Exec("DROP TABLE bookmarks_fts;").Error
				if err != nil {
					return err
				}

				err = tx.Exec("DROP TRIGGER bookmarks_ai;").Error
				if err != nil {
					return err
				}

				err = tx.Exec("DROP TRIGGER bookmarks_ad;").Error
				if err != nil {
					return err
				}

				err = tx.Exec("DROP TRIGGER bookmarks_au;").Error
				return err
			},
		},
	})
}
