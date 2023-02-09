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
		Privacy:     req.Privacy,
	}

	result := r.db.Create(bookmark)
	if result.Error != nil {
		return nil, result.Error
	}

	return bookmark, nil
}

func (r *BookmarkRepository) List(req BookmarkListRequest) (*BookmarkListResult, error) {
	query := r.db.Table("bookmarks")

	if req.Privacy != BookmarkPrivacyQueryAll {
		privacy := BookmarkPrivacyPublic
		if req.Privacy != "" {
			privacy = req.Privacy
		}
		query = query.Where("privacy = ?", privacy)
	}

	if req.Order != "" {
		query = query.Order(req.Order)
	}

	var count int64
	result := query.Count(&count)
	if result.Error != nil {
		return nil, result.Error
	}

	query.Offset(req.Offset)
	limit := 10
	query = query.Limit(limit)

	var bookmarks []Bookmark
	result = query.Find(&bookmarks)
	if result.Error != nil {
		return nil, result.Error
	}

	return &BookmarkListResult{
		Items:      bookmarks,
		Count:      count,
		HasPrev:    req.Offset > 0,
		PrevOffset: req.Offset - limit,
		HasNext:    int64(req.Offset+limit) < count,
		NextOffset: req.Offset + limit,
	}, nil
}
