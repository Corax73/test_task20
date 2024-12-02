package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"songLibrary/customLog"
	"songLibrary/models"
	"songLibrary/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Router struct {
	*mux.Router
}

var response map[string]interface{}

func (router *Router) Init() *Router {
	r := mux.NewRouter()
	r.HandleFunc("/songs/{id:[0-9]+}", router.getOneSongs).Methods("GET")
	r.HandleFunc("/songs/", router.CreateSong).Methods("POST")
	return &Router{r}
}

func (router *Router) initProcess(w http.ResponseWriter, r *http.Request, isPost bool) map[string]string {
	var resp map[string]string
	w.Header().Set("Content-Type", "application/json")
	if router.checkEnv() {
		router.consoleOutput(r)
	}
	if !isPost {
		resp = mux.Vars(r)
	} else {
		err := json.NewDecoder(r.Body).Decode(&resp)
		if err != nil {
			customLog.Logging(err)
		}
	}
	return resp
}

func (router *Router) checkEnv() bool {
	var resp bool
	envData := utils.GetConfFromEnvFile()
	if val, ok := envData["CONSOLE_OUT"]; ok && val == "true" {
		resp = true
	}
	return resp
}

func (router *Router) consoleOutput(r *http.Request) {
	fmt.Println(strings.Join([]string{time.Now().Format(time.RFC3339), r.Method, r.RequestURI, r.UserAgent()}, " "))
}

func (router *Router) getOneSongs(w http.ResponseWriter, r *http.Request) {
	params := router.initProcess(w, r, false)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		customLog.Logging(err)
	} else {
		songModel := (*&models.Song{}).Init()
		response = map[string]interface{}{"data": songModel.GetOneById(id)}
	}
	json.NewEncoder(w).Encode(response)
}

func (router *Router) CreateSong(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{}
	params := router.initProcess(w, r, true)
	if group, ok := params["group"]; ok && group != "" {
		groupModel := (*&models.Group{}).Init()
		songModel := (*&models.Song{}).Init()
		existGroup := groupModel.GetOneByTitle(map[string]string{"title": group})
		var groupId string
		if _, ok := existGroup["error"]; ok {
			result := groupModel.Create(map[string]string{"id": "", "title": group})
			if id, ok := result["id"]; !ok {
				response["data"] = "Error.Try again"
			} else {
				groupId = id
			}
		} else {
			groupId = strconv.FormatInt(existGroup["id"].(int64), 10)
		}
		if song, ok := params["song"]; ok && song != "" {
			result := songModel.Create(map[string]string{
				"id":          "",
				"title":       song,
				"group_id":    groupId,
				"link":        "",
				"releaseDate": "",
				"text":        "",
			})
			if id, ok := result["id"]; !ok {
				response["data"] = "Error.Try again"
			} else {
				response["id"] = id
			}
		} else {
			response["error"] = "Check parameters"
		}
	}
	json.NewEncoder(w).Encode(response)
}