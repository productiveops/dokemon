package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/productiveops/dokemon/pkg/server/model"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type FileSystemComposeLibraryStore struct {
	db *gorm.DB
	composeLibraryPath string
}

func NewFileSystemComposeLibraryStore(db *gorm.DB, composeLibraryPath string) *FileSystemComposeLibraryStore {
	return &FileSystemComposeLibraryStore{
		db: db,
		composeLibraryPath: composeLibraryPath,
	}
}

func (s *FileSystemComposeLibraryStore) Create(m *model.ComposeLibraryItem) error {
	p := filepath.Join(s.composeLibraryPath, m.ProjectName)

	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(p, 0755)
		if err != nil {
			return err
		}
	
		filename := filepath.Join(p, "compose.yaml")

		f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}

		f.WriteString(m.Definition)
		if err := f.Close(); err != nil {
			return err
		}

		return nil
	} else {
		return errors.New("Another project with this name already exists.")
	}
}

func (s *FileSystemComposeLibraryStore) Update(m *model.ComposeLibraryItemUpdate) error {
	composeProjectDirPath := filepath.Join(s.composeLibraryPath, m.ProjectName)
	_, err := os.ReadDir(filepath.Join(composeProjectDirPath))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("Project does not exist")
		} else {
			return err
		}
	}

	if m.ProjectName != m.NewProjectName {
		newComposeProjectDirPath := filepath.Join(s.composeLibraryPath, m.NewProjectName)
		_, err := os.Stat(filepath.Join(newComposeProjectDirPath))
		if err == nil {
			return errors.New("Another project with this name exist")
		}

		err = os.Rename(composeProjectDirPath, newComposeProjectDirPath)
		if err != nil {
			log.Error().Err(err).Msg("Error while renaming compose project directory")
			return err
		}

		composeProjectDirPath = newComposeProjectDirPath

		r := s.db.Exec("update node_compose_projects set library_project_name = ? where library_project_name = ?", m.NewProjectName, m.ProjectName)
		if r.Error != nil {
			log.Error().Err(err).Msg("Error while renaming compose project references")
			return err
		}
	}

	composeProjectFilePath := filepath.Join(composeProjectDirPath, "compose.yaml")
	_, err = os.Stat(composeProjectFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("Compose definition does not exist")
		} else {
			return err
		}
	}

	f, err := os.OpenFile(composeProjectFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	f.WriteString(m.Definition)
	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func (s *FileSystemComposeLibraryStore) GetByName(projectName string) (*model.ComposeLibraryItem, error) {
	composeProjectDirPath := filepath.Join(s.composeLibraryPath, projectName)
	_, err := os.ReadDir(filepath.Join(composeProjectDirPath))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.New("Project does not exist")
		} else {
			return nil, err
		}
	}

	composeProjectFilePath := filepath.Join(composeProjectDirPath, "compose.yaml")
	_, err = os.Stat(composeProjectFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.New("Compose definition does not exist")
		} else {
			return nil, err
		}	
	}

	definitionBytes, err := os.ReadFile(composeProjectFilePath)
	if err != nil {
		return nil, err
	}

	return &model.ComposeLibraryItem{ProjectName: projectName, Definition: string(definitionBytes)}, nil
}

func (s *FileSystemComposeLibraryStore) DeleteByName(projectName string) error {
	composeProjectDirPath := filepath.Join(s.composeLibraryPath, projectName)
	_, err := os.ReadDir(filepath.Join(composeProjectDirPath))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("Project does not exist")
		} else {
			return err
		}
	}

	var count int64
	err = s.db.Model(&model.NodeComposeProject{}).Where("library_project_name = ?", projectName).Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New(fmt.Sprintf("Definition is in use by %d projects and cannot be deleted", count))
	}

	err = os.RemoveAll(composeProjectDirPath)
	if err != nil {
		return err
	}

	return nil
}

func (s *FileSystemComposeLibraryStore) GetList(pageNo, pageSize uint) ([]model.ComposeLibraryItemHead, int64, error) {
	entries, err := os.ReadDir(s.composeLibraryPath)
	if err != nil {
		return nil, 0, err
	}

	composeItemHeads := make([]model.ComposeLibraryItemHead, len(entries))
	for i, entry := range entries {
		composeItemHeads[i] = model.ComposeLibraryItemHead{ProjectName: entry.Name()}
	}
	
	startIndex := (pageNo - 1) * pageSize
	endIndex := startIndex + pageSize
	if endIndex > uint(len(composeItemHeads)) {
		endIndex = uint(len(composeItemHeads))
	}
	return composeItemHeads[startIndex:endIndex], int64(len(entries)), nil
}

