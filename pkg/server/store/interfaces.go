package store

import (
	"time"

	"github.com/productiveops/dokemon/pkg/server/model"
)

type UserStore interface {
	Create(m *model.User) error
	Update(m *model.User) error
	GetById(id uint) (*model.User, error)
	GetByUserName(username string) (*model.User, error)
	GetList(pageNo, pageSize uint) ([]model.User, int64, error)
	DeleteById(id uint) error
	Exists(id uint) (bool, error)

	IsUniqueUserName(userName string) (bool, error)
	IsUniqueUserNameExcludeItself(userName string, id uint) (bool, error)

	Count() (int64, error)
}

type NodeStore interface {
	Create(m *model.Node) error
	Update(m *model.Node) error
	UpdateAgentVersion(id uint, version string) error
	UpdateLastPing(id uint, t time.Time) error
	UpdateContainerBaseUrl(id uint, url *string) error
	GetById(id uint) (*model.Node, error)
	GetList(pageNo, pageSize uint) ([]model.Node, int64, error)
	DeleteById(id uint) error
	Exists(id uint) (bool, error)

	IsUniqueName(name string) (bool, error)
	IsUniqueNameExcludeItself(name string, id uint) (bool, error)
}

type NodeComposeProjectStore interface {
	Create(m *model.NodeComposeProject) error
	Update(m *model.NodeComposeProject) error
	GetById(nodeId uint, id uint) (*model.NodeComposeProject, error)
	GetList(nodeId uint, pageNo, pageSize uint) ([]model.NodeComposeProject, int64, error)
	DeleteById(nodeId uint, id uint) error
	Exists(nodeId uint, id uint) (bool, error)

	IsUniqueName(nodeId uint, name string) (bool, error)
	IsUniqueNameExcludeItself(nodeId uint, name string, id uint) (bool, error)
}

type SettingStore interface {
	Create(m *model.Setting) error
	Update(m *model.Setting) error
	GetById(id string) (*model.Setting, error)
	DeleteById(id string) error
	Exists(id string) (bool, error)
}

type FileSystemComposeLibraryStore interface {
	Create(m *model.FileSystemComposeLibraryItem) error
	Update(m *model.FileSystemComposeLibraryItemUpdate) error
	GetByName(projectName string) (*model.FileSystemComposeLibraryItem, error)
	DeleteByName(projectName string) error
	GetList() ([]model.FileSystemComposeLibraryItemHead, int64, error)
	IsUniqueName(projectName string) (bool, error)
	IsUniqueNameExcludeItself(newProjectName string, existingProjectName string) (bool, error)
}

type ComposeLibraryStore interface {
	Create(m *model.ComposeLibraryItem) error
	Update(m *model.ComposeLibraryItem) error
	GetById(id uint) (*model.ComposeLibraryItem, error)
	GetList() ([]model.ComposeLibraryItem, int64, error)
	DeleteById(id uint) error
	Exists(id uint) (bool, error)

	IsUniqueName(name string) (bool, error)
	IsUniqueNameExcludeItself(name string, id uint) (bool, error)
}

type CredentialStore interface {
	Create(m *model.Credential) error
	Update(m *model.Credential) error
	GetById(id uint) (*model.Credential, error)
	GetList(pageNo, pageSize uint) ([]model.Credential, int64, error)
	DeleteById(id uint) error
	Exists(id uint) (bool, error)

	IsUniqueName(name string) (bool, error)
	IsUniqueNameExcludeItself(name string, id uint) (bool, error)
}

type EnvironmentStore interface {
	Create(m *model.Environment) error
	Update(m *model.Environment) error
	GetById(id uint) (*model.Environment, error)
	GetList(pageNo, pageSize uint) ([]model.Environment, int64, error)
	GetMap() (map[string]string, error)
	DeleteById(id uint) error
	Exists(id uint) (bool, error)

	IsUniqueName(name string) (bool, error)
	IsUniqueNameExcludeItself(name string, id uint) (bool, error)
}

type VariableStore interface {
	Create(m *model.Variable) error
	Update(m *model.Variable) error
	GetById(id uint) (*model.Variable, error)
	GetList(offset, limit uint) ([]model.Variable, int64, error)
	DeleteById(id uint) error
	Exists(id uint) (bool, error)

	IsUniqueName(name string) (bool, error)
	IsUniqueNameExcludeItself(name string, id uint) (bool, error)
}

type VariableValueStore interface {
	CreateOrUpdate(m *model.VariableValue) error
	Get(variableId, environmentId uint) (*model.VariableValue, error)
	GetMap(variableId uint) (map[string]string, error)
	GetMapByEnvironment(environmentId uint) (map[string]VariableValue, error)
}