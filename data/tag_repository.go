package data

import (
	"strings"

	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db}
}

func (r *TagRepository) GetByName(name string) (*Tag, error) {
	if name == "" {
		return nil, nil
	}

	var tag Tag
	result := r.db.Where("name = ?", strings.ToLower(name)).First(&tag)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &tag, nil
}

func (r *TagRepository) Upsert(tagNames []string) ([]Tag, error) {
	tags := []Tag{}

	for _, name := range tagNames {
		var tag Tag
		result := r.db.Where("name = ?", strings.ToLower(name)).First(&tag)
		if result.RowsAffected > 0 {
			tags = append(tags, tag)
		} else {
			tag = Tag{DisplayName: name}
			created := r.db.Create(&tag)
			if created.Error != nil {
				return nil, created.Error
			}
			tags = append(tags, tag)
		}
	}

	return tags, nil
}

func (r *TagRepository) Search(query string) ([]Tag, error) {
	var tags []Tag
	err := r.db.Table("tags_fts").Unscoped().Where("tags_fts MATCH ?", query).Order("rank").Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}
