package store

import (
	"dokemon/pkg/server/model"
	"errors"

	"gorm.io/gorm"
)

type SqlSettingStore struct {
	db *gorm.DB
}

func NewSqlSettingStore(db *gorm.DB) *SqlSettingStore {
	return &SqlSettingStore{
		db: db,
	}
}

func (s *SqlSettingStore) Create(m *model.Setting) error {
	return s.db.Create(m).Error
}

func (s *SqlSettingStore) Update(m *model.Setting) error {
	return s.db.Save(m).Error
}

func (s *SqlSettingStore) GetById(id string) (*model.Setting, error) {
	var m model.Setting

	if err := s.db.Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &m, nil
}

func (s *SqlSettingStore) Exists(id string) (bool, error) {
	var count int64

	if err := s.db.Model(&model.Setting{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *SqlSettingStore) DeleteById(id string) error {
	if err := s.db.Delete(&model.Setting{}, id).Error; err != nil {
		return err
	}

	return nil
}
