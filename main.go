package main

import (
	"fmt"
	"net/http"
	"songLibrary/customLog"
	extapi "songLibrary/extApi"
	"songLibrary/models"
	"songLibrary/repository"
	"songLibrary/router"
	"songLibrary/utils"
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
	mainPort := ":8080"
	fake3dApiPort := ":8082"
	envData := utils.GetConfFromEnvFile()
	if val, ok := envData["MAIN_PORT"]; ok {
		mainPort = utils.ConcatSlice([]string{":", val})
	}
	if val, ok := envData["FAKE_3D_API_PORT"]; ok {
		fake3dApiPort = utils.ConcatSlice([]string{":", val})
	}
	defer wg.Wait()
	wg.Add(1)
	go func(errChan chan<- error, handler http.Handler) {
		errChan <- http.ListenAndServe(mainPort, handler)
		defer wg.Done()
	}(errChan, r)

	go func() {
		check := true
		var invitationPrinted bool
		for check {
			if len(errChan) > 0 {
				fmt.Println(<-errChan)
				check = false
			} else {
				if !invitationPrinted {
					fmt.Println(strings.Join([]string{"started ",
						termlink.Link(
							utils.ConcatSlice([]string{"http://localhost", mainPort}),
							utils.ConcatSlice([]string{"http://localhost", mainPort}),
						)},
						" ",
					))
					invitationPrinted = true
				}
			}
		}
	}()

	rExt := (*&extapi.ExtRouter{}).Init()
	errChanExt := make(chan error, 1)
	defer close(errChanExt)
	defer wg.Wait()
	wg.Add(1)
	go func(errChanExt chan<- error, handler http.Handler) {
		errChanExt <- http.ListenAndServe(fake3dApiPort, handler)
		defer wg.Done()
	}(errChan, rExt)

	checkExt := true
	var invitationPrintedExt bool
	for checkExt {
		if len(errChanExt) > 0 {
			fmt.Println(<-errChanExt)
			checkExt = false
		} else {
			if !invitationPrintedExt {
				fmt.Println(strings.Join([]string{"started ",
					termlink.Link(
						utils.ConcatSlice([]string{"http://localhost", fake3dApiPort}),
						utils.ConcatSlice([]string{"http://localhost", fake3dApiPort}),
					)},
					" ",
				))
				invitationPrintedExt = true
			}
		}
	}
}
