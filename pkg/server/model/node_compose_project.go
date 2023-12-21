package model

type NodeComposeProject struct {
	Id                 uint
	NodeId             uint
	Node               Node
	EnvironmentId      *uint
	Environment        *Environment
	LibraryProjectId   *uint
	LibraryProject     *ComposeLibraryItem
	LibraryProjectName *string `gorm:"size:50"`
	ProjectName        string  `gorm:"size:50"`
	Type               string  `gorm:"size:20,default:''"` // github, local
	Url                *string `gorm:"size:255"`
	CredentialId       *uint
	Credential         *Credential
	Definition         *string // nil for github projects
}
