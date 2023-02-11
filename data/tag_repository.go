package data

import "gorm.io/gorm"

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db}
}

func (r *TagRepository) Upsert(tagNames []string) ([]Tag, error) {
	tags := []Tag{}

	for _, name := range tagNames {
		var tag Tag
		result := r.db.Where("name = ?", name).First(&tag)
		if result.RowsAffected > 0 {
			tags = append(tags, tag)
		} else {
			tag = Tag{Name: name}
			created := r.db.Create(&tag)
			if created.Error != nil {
				return nil, created.Error
			}
			tags = append(tags, tag)
		}
	}

	return tags, nil
}
