package model

type User struct {
	Id           uint
	UserName     string `gorm:"unique;size:255"`
	PasswordHash string `gorm:"size:255"`
}