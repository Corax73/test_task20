package main

import (
	"fmt"
	"songLibrary/customLog"
	"songLibrary/models"
	"songLibrary/repository"
)

type Data struct {
	Title string
}

func main() {
	customLog.LogInit("./logs/app.log")
	groupModel := (*&models.Group{}).Init()
	songModel := (*&models.Song{}).Init()
	modelsList := []*models.Model{
		groupModel.Model,
		songModel.Model,
	}
	var msg string
	if repository.Init(modelsList) {
		msg = "tables exist"
	} else {
		msg = "check the logs"
	}
	fmt.Println(msg)

	groupModel.Create(map[string]string{"id": "", "title": "test1"})
}
