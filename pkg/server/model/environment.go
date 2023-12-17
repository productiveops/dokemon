package model

type Environment struct {
	Id        			uint
	Name      			string      `gorm:"unique;size:50"`
}