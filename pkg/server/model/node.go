package model

import "time"

type Node struct {
	Id        			uint
	EnvironmentId      	*uint
	Environment        	*Environment
	Name      			string      `gorm:"unique;size:50"`
	AgentVersion		string		`gorm:"size:20"`
	TokenHash 			*string     `gorm:"unique;size:100"`
	LastPing  			*time.Time
	ContainerBaseUrl 	*string 	`gorm:"size:255"`
}
