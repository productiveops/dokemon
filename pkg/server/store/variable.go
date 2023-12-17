package store

import (
	"dokemon/pkg/server/model"
	"errors"

	"gorm.io/gorm"
)

type SqlVariableStore struct {
	db *gorm.DB
}

func NewSqlVariableStore(db *gorm.DB) *SqlVariableStore {
	return &SqlVariableStore{
		db: db,
	}
}

func (s *SqlVariableStore) Create(m *model.Variable) error {
	return s.db.Create(m).Error
}

func (s *SqlVariableStore) Update(m *model.Variable) error {
	return s.db.Save(m).Error
}

func (s *SqlVariableStore) GetById(id uint) (*model.Variable, error) {
	var m model.Variable

	if err := s.db.Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &m, nil
}

func (s *SqlVariableStore) Exists(id uint) (bool, error) {
	var count int64

	if err := s.db.Model(&model.Variable{}).Where("id = ?",  id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *SqlVariableStore) DeleteById(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error { 
		if err := tx.Where("variable_id = ?", id).Delete(&model.VariableValue{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&model.Variable{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *SqlVariableStore) GetList(pageNo, pageSize uint) ([]model.Variable, int64, error) {
	var (
		l []model.Variable
		count int64
	)

	s.db.Model(&l).Count(&count)
	s.db.Offset(int((pageNo - 1) * pageSize)).Limit(int(pageSize)).Order("id asc").Find(&l)

	return l, count, nil
}

func (s *SqlVariableStore) IsUniqueName(name string) (bool, error) {
	var count int64

	if err := s.db.Model(&model.Variable{}).Where("name = ? COLLATE NOCASE", name).Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil 
}

func (s *SqlVariableStore) IsUniqueNameExcludeItself(name string, id uint) (bool, error) {
	var count int64

	if err := s.db.Model(&model.Variable{}).Where("name = ? COLLATE NOCASE and id <> ?", name, id).Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil 
}
