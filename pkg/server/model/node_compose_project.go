package model

type NodeComposeProject struct {
	Id                 uint
	NodeId             uint
	Node               Node
	EnvironmentId      *uint
	Environment        *Environment
	LibraryProjectId   *uint
	LibraryProject     *ComposeLibraryItem
	LibraryProjectName string `gorm:"size:50"`
	ProjectName        string `gorm:"size:50"`
}
