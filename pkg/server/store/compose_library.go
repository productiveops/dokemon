package store

import (
	"errors"

	"github.com/productiveops/dokemon/pkg/server/model"

	"gorm.io/gorm"
)

type SqlComposeLibraryStore struct {
	db *gorm.DB
}

func NewSqlComposeLibraryStore(db *gorm.DB) *SqlComposeLibraryStore {
	return &SqlComposeLibraryStore{
		db: db,
	}
}

func (s *SqlComposeLibraryStore) Create(m *model.ComposeLibraryItem) error {
	return s.db.Create(m).Error
}

func (s *SqlComposeLibraryStore) Update(m *model.ComposeLibraryItem) error {
	return s.db.Save(m).Error
}

func (s *SqlComposeLibraryStore) GetById(id uint) (*model.ComposeLibraryItem, error) {
	var m model.ComposeLibraryItem

	if err := s.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &m, nil
}

func (s *SqlComposeLibraryStore) Exists(id uint) (bool, error) {
	var count int64

	if err := s.db.Model(&model.ComposeLibraryItem{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *SqlComposeLibraryStore) DeleteById(id uint) error {
	if err := s.db.Delete(&model.ComposeLibraryItem{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (s *SqlComposeLibraryStore) GetList() ([]model.ComposeLibraryItem, int64, error) {
	var (
		l []model.ComposeLibraryItem
		count int64
	)

	s.db.Model(&l).Count(&count)
	s.db.Order("project_name asc").Find(&l)

	return l, count, nil
}

func (s *SqlComposeLibraryStore) IsUniqueName(name string) (bool, error) {
	var count int64

	if err := s.db.Model(&model.ComposeLibraryItem{}).Where("project_name = ? COLLATE NOCASE", name).Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil 
}

func (s *SqlComposeLibraryStore) IsUniqueNameExcludeItself(name string, id uint) (bool, error) {
	var count int64

	if err := s.db.Model(&model.ComposeLibraryItem{}).Where("project_name = ? COLLATE NOCASE and id <> ?", name, id).Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil 
}
