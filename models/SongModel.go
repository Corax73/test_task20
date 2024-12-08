package models

type Song struct {
	*Model
}

func (song *Song) Init() *Song {
	model := Model{}
	model.SetTable("songs")
	model.Fields = map[string]string{"id": "", "group_id": "", "title": "", "release_date": "", "text": "", "link": ""}
	return &Song{&model}
}

func (song *Song) FireEvent() {

}