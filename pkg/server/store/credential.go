package store

import (
	"errors"

	"github.com/productiveops/dokemon/pkg/server/model"

	"gorm.io/gorm"
)

type SqlCredentialStore struct {
	db *gorm.DB
}

func NewSqlCredentialStore(db *gorm.DB) *SqlCredentialStore {
	return &SqlCredentialStore{
		db: db,
	}
}

func (s *SqlCredentialStore) Create(m *model.Credential) error {
	return s.db.Create(m).Error
}

func (s *SqlCredentialStore) Update(m *model.Credential) error {
	return s.db.Save(m).Error
}

func (s *SqlCredentialStore) GetById(id uint) (*model.Credential, error) {
	var m model.Credential

	if err := s.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &m, nil
}

func (s *SqlCredentialStore) Exists(id uint) (bool, error) {
	var count int64

	if err := s.db.Model(&model.Credential{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *SqlCredentialStore) IsInUse(id uint) (bool, error) {
	var ncp_ref_count, cli_ref_count int64

	if err := s.db.Model(&model.NodeComposeProject{}).Where("credential_id = ?", id).Count(&ncp_ref_count).Error; err != nil {
		return false, err
	}

	if err := s.db.Model(&model.ComposeLibraryItem{}).Where("credential_id = ?", id).Count(&cli_ref_count).Error; err != nil {
		return false, err
	}

	return (ncp_ref_count + cli_ref_count) > 0, nil
}

func (s *SqlCredentialStore) DeleteById(id uint) error {
	inUse, err := s.IsInUse(id)
	if err != nil {
		return err
	}

	if inUse {
		return errors.New("Credentials are in use and cannot be deleted")
	}

	if err := s.db.Delete(&model.Credential{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (s *SqlCredentialStore) GetList(pageNo, pageSize uint) ([]model.Credential, int64, error) {
	var (
		l []model.Credential
		count int64
	)

	s.db.Model(&l).Count(&count)
	s.db.Offset(int((pageNo - 1) * pageSize)).Limit(int(pageSize)).Order("name asc").Find(&l)

	return l, count, nil
}

func (s *SqlCredentialStore) IsUniqueName(name string) (bool, error) {
	var count int64

	if err := s.db.Model(&model.Credential{}).Where("name = ? COLLATE NOCASE", name).Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil 
}

func (s *SqlCredentialStore) IsUniqueNameExcludeItself(name string, id uint) (bool, error) {
	var count int64

	if err := s.db.Model(&model.Credential{}).Where("name = ? COLLATE NOCASE and id <> ?", name, id).Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil 
}
