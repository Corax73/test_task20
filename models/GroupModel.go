package models

import (
	"songLibrary/customDb"
	"songLibrary/utils"
)

type Group struct {
	*Model
}

func (group *Group) Init() *Group {
	model := Model{}
	model.SetTable("groups")
	model.Fields = map[string]string{"id": "", "title": ""}
	return &Group{&model}
}

func (group *Group) Create(fields map[string]string) bool {
	var resp bool
	if utils.CompareMapsByStringKeys(group.Fields, fields) {
		group.Fields = fields
		db := customDb.GetConnect()
		defer customDb.CloseConnect(db)
		resp = group.Save()
	}
	return resp
}
