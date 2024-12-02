package main

import (
	"fmt"
	"net/http"
	"songLibrary/customLog"
	"songLibrary/models"
	"songLibrary/repository"
	"songLibrary/router"
	"strings"
	"sync"

	"github.com/savioxavier/termlink"
)

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
	} else {
		msg = "check the logs"
	}
	fmt.Println(msg)
	r := (*&router.Router{}).Init()
	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	defer close(errChan)
	defer wg.Wait()
	wg.Add(1)
	go func(errChan chan<- error, handler http.Handler) {
		errChan <- http.ListenAndServe(":8000", handler)
		defer wg.Done()
	}(errChan, r)

	check := true
	var invitationPrinted bool
	for check {
		if len(errChan) > 0 {
			fmt.Println(<-errChan)
			check = false
		} else {
			if !invitationPrinted {
				fmt.Println(strings.Join([]string{"started ", termlink.Link("http://localhost:8000", "http://localhost:8000")}, " "))
				invitationPrinted = true
			}
		}
	}
}
