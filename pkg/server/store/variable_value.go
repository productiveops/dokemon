package store

import (
	"errors"
	"strconv"

	"github.com/productiveops/dokemon/pkg/crypto/ske"
	"github.com/productiveops/dokemon/pkg/server/model"

	"gorm.io/gorm"
)

type SqlVariableValueStore struct {
	db *gorm.DB
}

func NewSqlVariableValueStore(db *gorm.DB) *SqlVariableValueStore {
	return &SqlVariableValueStore{
		db: db,
	}
}

func (s *SqlVariableValueStore) CreateOrUpdate(m *model.VariableValue) error {
	var v model.VariableValue

	if err := s.db.Where("variable_id = ? and environment_id = ?", m.VariableId, m.EnvironmentId).First(&v).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.db.Create(m).Error
		}
		return err
	}

	v.Value = m.Value

	return s.db.Save(v).Error
}

func (s *SqlVariableValueStore) Get(variableId, environmentId uint) (*model.VariableValue, error) {
	var m model.VariableValue

	if err := s.db.Where("variable_id = ? and environment_id = ?", variableId, environmentId).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}  else {
			return nil, err
		}
	}

	return &m, nil
}

type envToValue struct {
	Id 		uint
	Value 	*string
}

func (s *SqlVariableValueStore) GetMap(variableId uint) (map[string]string, error) {
	ret := make(map[string]string)

	var v []envToValue

	if err := s.db.Model(&model.Environment{}).
				Select("environments.id, variable_values.value").
				Joins("left join variable_values on environments.id = variable_values.environment_id and variable_values.variable_id = ?", variableId).
				Scan(&v).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	for _, row := range v {
		if row.Value == nil {
			ret[strconv.Itoa(int(row.Id))] = ""
		} else {
			decryptedValue, err := ske.Decrypt(*row.Value)
			if err != nil {
				return nil, err
			}
			ret[strconv.Itoa(int(row.Id))] = decryptedValue
		}
	}

	return ret, nil
}

type variable struct {
	Name 	string
	Value 	*string
	IsSecret bool
}

type VariableValue struct {
	Value 	*string
	IsSecret bool
}

func (s *SqlVariableValueStore) GetMapByEnvironment(environmentId uint) (map[string]VariableValue, error) {
	ret := make(map[string]VariableValue)

	var variables []variable

	if err := s.db.Model(&model.Variable{}).
				Select("variables.name, variable_values.value, variables.is_secret").
				Joins("left join variable_values on variable_values.variable_id = variables.id and variable_values.environment_id = ?", environmentId).
				Scan(&variables).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	for _, variable := range variables {
		value := ""

		if variable.Value != nil {
			decryptedValue, err := ske.Decrypt(*variable.Value)
			if err != nil {
				return nil, err
			}
			value = decryptedValue
		}

		ret[variable.Name] = VariableValue{
			Value: &value,
			IsSecret: variable.IsSecret,
		}
	}

	return ret, nil
}
