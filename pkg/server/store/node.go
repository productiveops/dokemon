package store

import (
	"dokemon/pkg/server/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

type SqlNodeStore struct {
	db *gorm.DB
}

func NewSqlNodeStore(db *gorm.DB) *SqlNodeStore {
	return &SqlNodeStore{
		db: db,
	}
}

func (s *SqlNodeStore) Create(m *model.Node) error {
	return s.db.Create(m).Error
}

func (s *SqlNodeStore) Update(m *model.Node) error {
	return s.db.Save(m).Error
}

func (s *SqlNodeStore) UpdateAgentVersion(id uint, version string) error {
	db := s.db.Model(&model.Node{}).Where("id = ?", id).Update("agent_version", version)

	return db.Error
}

func (s *SqlNodeStore) UpdateLastPing(id uint, t time.Time) error {
	db := s.db.Model(&model.Node{}).Where("id = ?", id).Update("last_ping", t)

	return db.Error
}

func (s *SqlNodeStore) UpdateContainerBaseUrl(id uint, url *string) error {
	db := s.db.Model(&model.Node{}).Where("id = ?", id).Update("container_base_url", url)

	return db.Error
}

func (s *SqlNodeStore) GetById(id uint) (*model.Node, error) {
	var m model.Node

	if err := s.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &m, nil
}

func (s *SqlNodeStore) Exists(id uint) (bool, error) {
	var count int64

	if err := s.db.Model(&model.Node{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *SqlNodeStore) DeleteById(id uint) error {
	if err := s.db.Delete(&model.Node{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (s *SqlNodeStore) GetList(pageNo, pageSize uint) ([]model.Node, int64, error) {
	var (
		l []model.Node
		count int64
	)

	s.db.Model(&l).Count(&count)
	s.db.Offset(int((pageNo - 1) * pageSize)).Limit(int(pageSize)).Order("nodes.name asc").Model(&model.Node{}).Joins("Environment").Find(&l)

	return l, count, nil
}

func (s *SqlNodeStore) IsUniqueName(name string) (bool, error) {
	var count int64

	if err := s.db.Model(&model.Node{}).Where("name = ? COLLATE NOCASE", name).Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil 
}

func (s *SqlNodeStore) IsUniqueNameExcludeItself(name string, id uint) (bool, error) {
	var count int64

	if err := s.db.Model(&model.Node{}).Where("name = ? COLLATE NOCASE and id <> ?", name, id).Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil 
}
