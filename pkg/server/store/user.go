package store

import (
	"dokemon/pkg/server/model"
	"errors"

	"gorm.io/gorm"
)

type SqlUserStore struct {
	db *gorm.DB
}

func NewSqlUserStore(db *gorm.DB) *SqlUserStore {
	return &SqlUserStore{
		db: db,
	}
}

func (s *SqlUserStore) Create(m *model.User) error {
	return s.db.Create(m).Error
}

func (s *SqlUserStore) Update(m *model.User) error {
	return s.db.Save(m).Error
}

func (s *SqlUserStore) GetById(id uint) (*model.User, error) {
	var m model.User

	if err := s.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &m, nil
}

func (s *SqlUserStore) GetByUserName(username string) (*model.User, error) {
	var m model.User

	if err := s.db.Where("user_name = ?", username).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &m, nil
}

func (s *SqlUserStore) Exists(id uint) (bool, error) {
	var count int64

	if err := s.db.Model(&model.User{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *SqlUserStore) DeleteById(id uint) error {
	if err := s.db.Delete(&model.User{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (s *SqlUserStore) GetList(pageNo, pageSize uint) ([]model.User, int64, error) {
	var (
		l []model.User
		count int64
	)

	s.db.Model(&l).Count(&count)
	s.db.Offset(int((pageNo - 1) * pageSize)).Limit(int(pageSize)).Order("user_name asc").Find(&l)

	return l, count, nil
}

func (s *SqlUserStore) IsUniqueUserName(userName string) (bool, error) {
	var count int64

	if err := s.db.Model(&model.User{}).Where("user_name = ? COLLATE NOCASE", userName).Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil 
}

func (s *SqlUserStore) IsUniqueUserNameExcludeItself(userName string, id uint) (bool, error) {
	var count int64

	if err := s.db.Model(&model.User{}).Where("user_name = ? COLLATE NOCASE and id <> ?", userName, id).Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil 
}

func (s *SqlUserStore) IsUniqueEmail(email string) (bool, error) {
	var count int64

	if err := s.db.Model(&model.User{}).Where("email = ? COLLATE NOCASE", email).Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil 
}

func (s *SqlUserStore) IsUniqueEmailExcludeItself(email string, id uint) (bool, error) {
	var count int64

	if err := s.db.Model(&model.User{}).Where("email = ? COLLATE NOCASE and id <> ?", email, id).Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil 
}

func (s *SqlUserStore) Count() (int64, error) {
	var count int64

	if err := s.db.Model(&model.User{}).Count(&count).Error; err != nil {
		return -1, err
	}

	return count, nil
}
