package main

import (
	"fmt"
	"songLibrary/customLog"
	"songLibrary/models"
	"songLibrary/repository"
	"strconv"
	"time"
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
	rep := repository.Repository{}
	if rep.Init(modelsList) {
		msg = "tables exist"
		fmt.Println(groupModel.Create(map[string]string{"id": "", "title": strconv.Itoa(time.Now().Second())}))
	} else {
		msg = "check the logs"
	}
	fmt.Println(msg)
}
