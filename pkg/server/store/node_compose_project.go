package store

import (
	"errors"

	"github.com/productiveops/dokemon/pkg/server/model"

	"gorm.io/gorm"
)

type SqlNodeComposeProjectStore struct {
	db *gorm.DB
	composeLibraryPath string
}

func NewSqlNodeComposeProjectStore(db *gorm.DB, composeLibraryPath string) *SqlNodeComposeProjectStore {
	return &SqlNodeComposeProjectStore{
		db: db,
		composeLibraryPath: composeLibraryPath,
	}
}

func (s *SqlNodeComposeProjectStore) Create(m *model.NodeComposeProject) error {
	return s.db.Create(m).Error
}

func (s *SqlNodeComposeProjectStore) Update(m *model.NodeComposeProject) error {
	return s.db.Save(m).Error
}

func (s *SqlNodeComposeProjectStore) GetById(nodeId uint, id uint) (*model.NodeComposeProject, error) {
	var m model.NodeComposeProject

	if err := s.db.Where("node_id = ?", nodeId).First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &m, nil
}

func (s *SqlNodeComposeProjectStore) Exists(nodeId uint, id uint) (bool, error) {
	var count int64

	if err := s.db.Model(&model.NodeComposeProject{}).Where("node_id = ? and id = ?", nodeId, id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *SqlNodeComposeProjectStore) DeleteById(nodeId uint, id uint) error {
	if err := s.db.Where("node_id = ?", nodeId).Delete(&model.NodeComposeProject{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (s *SqlNodeComposeProjectStore) GetList(nodeId uint, pageNo, pageSize uint) ([]model.NodeComposeProject, int64, error) {
	var (
		l []model.NodeComposeProject
		count int64
	)

	s.db.Model(&l).Where("node_id = ?", nodeId).Count(&count)
	s.db.Where("node_id = ?", nodeId).Offset(int((pageNo - 1) * pageSize)).Limit(int(pageSize)).Order("project_name asc").Find(&l)

	return l, count, nil
}

func (s *SqlNodeComposeProjectStore) IsUniqueName(nodeId uint, name string) (bool, error) {
	var count int64

	if err := s.db.Model(&model.NodeComposeProject{}).Where("node_id = ? COLLATE NOCASE and project_name = ?", nodeId, name).Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil 
}

func (s *SqlNodeComposeProjectStore) IsUniqueNameExcludeItself(nodeId uint, name string, id uint) (bool, error) {
	var count int64

	if err := s.db.Model(&model.NodeComposeProject{}).
					Where("node_id = ? and project_name = ? COLLATE NOCASE and id <> ?", nodeId, name, id).
					Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil 
}

func (s *SqlNodeComposeProjectStore) UpdateOldVersionRecords() error {
	var (
		l []model.NodeComposeProject
	)

	if err := s.db.Where("type is null or type = ''").Find(&l).Error; err != nil {
		return err
	}

	for _, item := range l {
		if item.LibraryProjectId != nil { // GitHub project
			composeLibraryStore := NewSqlComposeLibraryStore(s.db)
			gitHubLibraryProject, err := composeLibraryStore.GetById(*item.LibraryProjectId)
			if err != nil {
				return err
			}
			item.CredentialId = gitHubLibraryProject.CredentialId
			item.Type = gitHubLibraryProject.Type
			item.Url = &gitHubLibraryProject.Url
		} else { // FileSystem project
			fileSystemComposeLibraryStore := NewLocalFileSystemComposeLibraryStore(s.db, s.composeLibraryPath)
			fileSystemLibraryProject, err := fileSystemComposeLibraryStore.GetByName(item.LibraryProjectName)
			if err != nil {
				return err
			}
			item.Type = "local"
			item.Definition = &fileSystemLibraryProject.Definition
		}

		if err := s.db.Save(item).Error; err != nil {
			return err
		}
	}

	return nil
}