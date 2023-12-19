package store

import (
	"errors"
	"fmt"

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

func (s *SqlComposeLibraryStore) GetByName(projectName string) (*model.ComposeLibraryItem, error) {
	var m model.ComposeLibraryItem

	if err := s.db.Where("project_name = ? COLLATE NOCASE", projectName).First(&m).Error; err != nil {
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
	m, err := s.GetById(id)
	if err != nil {
		return err
	}

	var count int64
	err = s.db.Model(&model.NodeComposeProject{}).Where("library_project_name = ?", m.ProjectName).Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New(fmt.Sprintf("Definition is in use by %d projects and cannot be deleted", count))
	}

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
