package store

import (
	"errors"

	"github.com/productiveops/dokemon/pkg/server/model"

	"gorm.io/gorm"
)

type SqlNodeComposeProjectVariableStore struct {
	db *gorm.DB
}

func NewSqlNodeComposeProjectVariableStore(db *gorm.DB) *SqlNodeComposeProjectVariableStore {
	return &SqlNodeComposeProjectVariableStore{
		db: db,
	}
}

func (s *SqlNodeComposeProjectVariableStore) Create(m *model.NodeComposeProjectVariable) error {
	return s.db.Create(m).Error
}

func (s *SqlNodeComposeProjectVariableStore) Update(m *model.NodeComposeProjectVariable) error {
	return s.db.Save(m).Error
}

func (s *SqlNodeComposeProjectVariableStore) GetById(nodeComposeProjectId uint, id uint) (*model.NodeComposeProjectVariable, error) {
	var m model.NodeComposeProjectVariable

	if err := s.db.Where("node_compose_project_id = ?", nodeComposeProjectId).First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &m, nil
}

func (s *SqlNodeComposeProjectVariableStore) Exists(nodeComposeProjectId uint, id uint) (bool, error) {
	var count int64

	if err := s.db.Model(&model.NodeComposeProjectVariable{}).Where("node_compose_project_id = ? and id = ?", nodeComposeProjectId, id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *SqlNodeComposeProjectVariableStore) DeleteById(nodeComposeProjectId uint, id uint) error {
	if err := s.db.Where("node_compose_project_id = ?", nodeComposeProjectId).Delete(&model.NodeComposeProjectVariable{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (s *SqlNodeComposeProjectVariableStore) GetList(nodeComposeProjectId uint, pageNo, pageSize uint) ([]model.NodeComposeProjectVariable, int64, error) {
	var (
		l []model.NodeComposeProjectVariable
		count int64
	)

	s.db.Model(&l).Where("node_compose_project_id = ?", nodeComposeProjectId).Count(&count)
	s.db.Where("node_compose_project_id = ?", nodeComposeProjectId).Offset(int((pageNo - 1) * pageSize)).Limit(int(pageSize)).Order("name asc").Find(&l)

	return l, count, nil
}

func (s *SqlNodeComposeProjectVariableStore) IsUniqueName(nodeComposeProjectId uint, name string) (bool, error) {
	var count int64

	if err := s.db.Model(&model.NodeComposeProjectVariable{}).
					Where("node_compose_project_id = ? COLLATE NOCASE and name = ?", nodeComposeProjectId, name).
					Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil 
}

func (s *SqlNodeComposeProjectVariableStore) IsUniqueNameExcludeItself(nodeComposeProjectId uint, name string, id uint) (bool, error) {
	var count int64

	if err := s.db.Model(&model.NodeComposeProjectVariable{}).
					Where("node_compose_project_id = ? and name = ? COLLATE NOCASE and id <> ?", nodeComposeProjectId, name, id).
					Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil 
}
