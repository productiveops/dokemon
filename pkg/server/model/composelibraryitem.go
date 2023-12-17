package model

// This is not a DB model

type ComposeLibraryItemHead struct {
	ProjectName string
}

type ComposeLibraryItem struct {
	ProjectName string
	Definition 	string
}

type ComposeLibraryItemUpdate struct {
	ProjectName string
	NewProjectName string
	Definition 	string
}
