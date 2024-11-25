package main

import (
	"fmt"
	"songLibrary/customDb"
	"songLibrary/customLog"
	"songLibrary/models"
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
	if customDb.Init(modelsList) {
		msg = "tables exist"
	} else {
		msg = "check the logs"
	}
	fmt.Println(msg)
}
