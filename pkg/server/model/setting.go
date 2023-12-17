package model

type Setting struct {
	Id    string `gorm:"size:100"`
	Value string
}