package data

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type BookmarkRepository struct {
	db *gorm.DB
}

func NewBookmarkRepository(db *gorm.DB) *BookmarkRepository {
	return &BookmarkRepository{db}
}

func (r *BookmarkRepository) Get(id uint) (*Bookmark, error) {
	var bookmark Bookmark
	result := r.db.Preload("Tags").First(&bookmark, id)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &bookmark, nil
}

func (r *BookmarkRepository) Create(form BookmarkForm) (*Bookmark, error) {
	var bookmark *Bookmark
	err := r.db.Transaction(func(tx *gorm.DB) error {
		tagRepo := NewTagRepository(tx)
		tagNames := parseTags(form.Tags)
		tags, err := tagRepo.Upsert(tagNames)
		if err != nil {
			return err
		}

		bookmark = &Bookmark{
			URL:         form.URL,
			Title:       form.Title,
			Description: form.Description,
			Privacy:     publicToPrivacy(form.Public),
			Tags:        tags,
		}

		result := tx.Create(bookmark)
		if result.Error != nil {
			bookmark = nil
			return result.Error
		}

		return nil
	})

	return bookmark, err
}

func (r *BookmarkRepository) List(req BookmarkListRequest) (*BookmarkListResult, error) {
	query := r.db.Model(&Bookmark{}).Preload("Tags")

	if req.Privacy != BookmarkPrivacyQueryAll {
		privacy := BookmarkPrivacyPublic
		if req.Privacy != "" {
			privacy = req.Privacy
		}
		query = query.Where("privacy = ?", privacy)
	}

	if req.TagID != 0 {
		query = query.Joins("inner join bookmark_tags bt on bt.bookmark_id = bookmarks.id").
			Where("bt.tag_id = ?", req.TagID)
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
		Items:   bookmarks,
		Count:   count,
		HasPrev: req.Offset > 0,
		PrevURL: fmt.Sprintf("%soffset=%d", req.PaginationPathPrefix, req.Offset-limit),
		HasNext: int64(req.Offset+limit) < count,
		NextURL: fmt.Sprintf("%soffset=%d", req.PaginationPathPrefix, req.Offset+limit),
	}, nil
}

func (r *BookmarkRepository) Delete(id uint) error {
	result := r.db.Delete(&Bookmark{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("bookmark with id %d not found", id)
	}

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *BookmarkRepository) Update(id uint, form BookmarkForm) error {
	bookmark, err := r.Get(id)
	if err != nil {
		return err
	}
	if bookmark == nil {
		return fmt.Errorf("bookmark with id %d not found", id)
	}

	var result *gorm.DB
	err = r.db.Transaction(func(tx *gorm.DB) error {
		bookmark.URL = form.URL
		bookmark.Title = form.Title
		bookmark.Description = form.Description
		bookmark.Privacy = publicToPrivacy(form.Public)

		result = tx.Save(bookmark)
		if result.Error != nil {
			return result.Error
		}

		tagRepo := NewTagRepository(tx)
		newTagNames := parseTags(form.Tags)
		updatedTags, err := tagRepo.Upsert(newTagNames)
		if err != nil {
			return err
		}
		err = tx.Model(&bookmark).Association("Tags").Replace(updatedTags)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func publicToPrivacy(public bool) BookmarkPrivacy {
	if public {
		return BookmarkPrivacyPublic
	}
	return BookmarkPrivacyPrivate
}

func parseTags(tagsString string) []string {
	tags := []string{}
	for _, name := range strings.Split(tagsString, ",") {
		trimmed := strings.TrimSpace(name)
		if trimmed != "" {
			tags = append(tags, strings.TrimSpace(name))
		}
	}
	return tags
}
