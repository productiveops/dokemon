package model

type Variable struct {
	Id       uint
	Name     string `gorm:"unique;size:100"`
	IsSecret bool
}