package main

import (
	"fmt"
	"songLibrary/customLog"
	"songLibrary/models"
)

type Data struct {
	Title string
}

func main() {
	customLog.LogInit("./logs/app.log")
	model := models.Model{}
	model.SetTable("test")
	fmt.Println(model.CheckModelTable())
	fmt.Println(model.RunTableMigration())
}
