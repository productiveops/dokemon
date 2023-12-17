package model

type VariableValue struct {
	Id            uint
	VariableId    uint
	Variable      Variable
	EnvironmentId uint
	Environment   Environment
	Value         string
}