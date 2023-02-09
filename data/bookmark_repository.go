package data

import "gorm.io/gorm"

type BookmarkRepository struct {
	db *gorm.DB
}

func NewBookmarkRepository(db *gorm.DB) *BookmarkRepository {
	return &BookmarkRepository{db}
}

func (r *BookmarkRepository) Create(req BookmarkCreate) (*Bookmark, error) {
	bookmark := &Bookmark{
		URL:         req.URL,
		Title:       req.Title,
		Description: req.Description,
	}

	result := r.db.Create(bookmark)
	if result.Error != nil {
		return nil, result.Error
	}

	return bookmark, nil
}
