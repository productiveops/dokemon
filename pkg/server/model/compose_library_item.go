package model

// FileSystem: This is not a DB model
type FileSystemComposeLibraryItemHead struct {
	ProjectName string
}

type FileSystemComposeLibraryItem struct {
	ProjectName string
	Definition 	string
}

type FileSystemComposeLibraryItemUpdate struct {
	ProjectName string
	NewProjectName string
	Definition 	string
}

// Remote: This is a DB model
type ComposeLibraryItem struct {
	Id				uint
	CredentialId 	*uint
	Credential 		*Credential
	ProjectName 	string
	Type 			string	// github
	Url 			string
}
