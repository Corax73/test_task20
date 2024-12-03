package models

type Models interface {
	Table() string
}

type HasEvent interface {
	FireEvent()
}