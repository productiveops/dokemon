package model

type NodeComposeProjectVariable struct {
	Id                   uint
	NodeComposeProjectId uint
	NodeComposeProject   NodeComposeProject
	Name                 string `gorm:"size:100"`
	IsSecret             bool
	Value                string
}