package model

// Local: This is not a DB model
type LocalComposeLibraryItemHead struct {
	ProjectName string
}

type LocalComposeLibraryItem struct {
	ProjectName string
	Definition 	string
}

type LocalComposeLibraryItemUpdate struct {
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
