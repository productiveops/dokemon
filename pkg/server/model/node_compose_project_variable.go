package model

type NodeComposeProjectVariable struct {
	Id                   uint
	NodeComposeProjectId uint
	NodeComposeProject   uint
	Name                 string `gorm:"unique;size:100"`
	IsSecret             bool
	Value                string
}