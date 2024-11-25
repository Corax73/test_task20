package models

type Group struct {
	*Model
}

func (group *Group) Init() *Group {
	model := Model{}
	model.SetTable("groups")
	model.SetFields(map[string]string{"id": "", "title": ""})
	return &Group{&model}
}
