package data

import (
	"errors"

	"gorm.io/gorm"

	"github.com/chdorner/submarine/util"
)

type SettingsRepository struct {
	db *gorm.DB
}

func NewSettingsRepository(db *gorm.DB) *SettingsRepository {
	return &SettingsRepository{db}
}

func (r *SettingsRepository) IsInitialized() bool {
	var count int64
	r.db.Model(&Settings{}).Count(&count)
	return count == 1
}

func (r *SettingsRepository) Get() (*Settings, error) {
	var settings Settings
	result := r.db.First(&settings)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &settings, nil
}

func (r *SettingsRepository) Upsert(value SettingsUpsert) error {
	existing, err := r.Get()
	if err != nil {
		return err
	}

	if existing == nil {
		return r.insert(&value)
	} else {
		return r.update(existing, &value)
	}
}

func (r *SettingsRepository) insert(req *SettingsUpsert) error {
	if req.Password == "" {
		return errors.New("password is empty")
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return err
	}
	value := Settings{
		Password: hashedPassword,
	}

	result := r.db.Create(&value)
	return result.Error
}

func (r *SettingsRepository) update(existing *Settings, req *SettingsUpsert) error {
	if req.Password != "" {
		hashedPassword, err := util.HashPassword(req.Password)
		if err != nil {
			return err
		}
		existing.Password = hashedPassword
	}

	result := r.db.Save(existing)
	return result.Error
}
