package store

import (
	"errors"
	"strconv"

	"github.com/productiveops/dokemon/pkg/server/model"

	"gorm.io/gorm"
)

type SqlEnvironmentStore struct {
	db *gorm.DB
}

func NewSqlEnvironmentStore(db *gorm.DB) *SqlEnvironmentStore {
	return &SqlEnvironmentStore{
		db: db,
	}
}

func (s *SqlEnvironmentStore) Create(m *model.Environment) error {
	return s.db.Create(m).Error
}

func (s *SqlEnvironmentStore) Update(m *model.Environment) error {
	return s.db.Save(m).Error
}

func (s *SqlEnvironmentStore) GetById(id uint) (*model.Environment, error) {
	var m model.Environment

	if err := s.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &m, nil
}

func (s *SqlEnvironmentStore) Exists(id uint) (bool, error) {
	var count int64

	if err := s.db.Model(&model.Environment{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *SqlEnvironmentStore) IsInUse(id uint) (bool, error) {
	var node_ref_count int64

	if err := s.db.Model(&model.Node{}).Where("environment_id = ?", id).Count(&node_ref_count).Error; err != nil {
		return false, err
	}

	return node_ref_count > 0, nil
}

func (s *SqlEnvironmentStore) DeleteById(id uint) error {
	inUse, err := s.IsInUse(id)
	if err != nil {
		return err
	}

	if inUse {
		return errors.New("Environment is in use and cannot be deleted")
	}

	if err := s.db.Delete(&model.Environment{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (s *SqlEnvironmentStore) GetList(pageNo, pageSize uint) ([]model.Environment, int64, error) {
	var (
		l []model.Environment
		count int64
	)

	s.db.Model(&l).Count(&count)
	s.db.Offset(int((pageNo - 1) * pageSize)).Limit(int(pageSize)).Order("name asc").Find(&l)

	return l, count, nil
}

func (s *SqlEnvironmentStore) IsUniqueName(name string) (bool, error) {
	var count int64

	if err := s.db.Model(&model.Environment{}).Where("name = ? COLLATE NOCASE", name).Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil 
}

func (s *SqlEnvironmentStore) IsUniqueNameExcludeItself(name string, id uint) (bool, error) {
	var count int64

	if err := s.db.Model(&model.Environment{}).Where("name = ? COLLATE NOCASE and id <> ?", name, id).Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil 
}

func (s *SqlEnvironmentStore) GetMap() (map[string]string, error) {
	l := []model.Environment{}
	s.db.Find(&l)

	m := make(map[string]string)

	for _, r := range l {
		m[strconv.Itoa(int(r.Id))] = r.Name
	}

	return m, nil
}