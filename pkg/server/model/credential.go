package model

type Credential struct {
	Id       uint
	Name     string  `gorm:"unique;size:50"`
	Service  *string `gorm:"size:50"`  // github
	Type     string  `gorm:"size:50"`  // pat = Personal Access Token
	UserName *string `gorm:"size:100"` // Username or any other identifier for the service
	Secret   string  // Token, password, etc.
}