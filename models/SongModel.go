package models

type Song struct {
	*Model
}

func (group *Song) Init() *Song {
	model := Model{}
	model.SetTable("songs")
	model.SetFields(map[string]string{"id": "", "group_id": "", "title": "", "releaseDate": "", "text": "", "link": ""})
	return &Song{&model}
}
