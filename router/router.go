package router

import (
	"encoding/json"
	"errors"
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
	r.HandleFunc("/songs/", router.getSongs).Methods("GET")
	r.HandleFunc("/songs/{id:[0-9]+}", router.getOneSongs).Methods("GET")
	r.HandleFunc("/songs/", router.CreateSong).Methods("POST")
	r.HandleFunc("/songs/{id:[0-9]+}", router.deleteSong).Methods("DELETE")
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
		sort := r.URL.Query().Get("sort")
		if sort != "" {
			var order string
			splits := strings.Split(sort, "--")
			if len(splits) > 1 {
				requestField, requestOrder := splits[0], splits[1]
				if requestOrder != "desc" && requestOrder != "asc" {
					order = "desc"
				} else {
					order = requestOrder
				}
				resp["order"] = order
				resp["orderBy"] = requestField
			}
		}
		filter := r.URL.Query().Get("filter")
		if filter != "" {
			splits := strings.Split(filter, "--")
			if len(splits) > 1 {
				requestField, requestValue := splits[0], splits[1]
				if requestValue != "" {
					resp["filterBy"] = requestField
					resp["filterVal"] = requestValue
				}
			}
		}
		limit := r.URL.Query().Get("limit")
		if limit != "" {
			resp["limit"] = limit
		}
		offset := r.URL.Query().Get("offset")
		if offset != "" {
			resp["offset"] = offset
		}
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

func (router *Router) getSongs(w http.ResponseWriter, r *http.Request) {
	params := router.initProcess(w, r, false)
	songModel := (*&models.Song{}).Init()
	response = map[string]interface{}{"data": songModel.GetList(params)}
	json.NewEncoder(w).Encode(response)
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
				"id":           "",
				"title":        song,
				"group_id":     groupId,
				"link":         "",
				"release_date": "",
				"text":         "",
			})
			if id, ok := result["id"]; !ok {
				response["data"] = "Error.Try again"
			} else {
				response["id"] = id
				if songModel.CheckInterface(songModel) {
					resp, err := http.Get(utils.ConcatSlice([]string{"http://localhost:8082/info/", group, "/", song}))
					if err != nil {
						customLog.Logging(err)
					} else {
						defer resp.Body.Close()
						var data map[string]string
						if err := json.NewDecoder(resp.Body).Decode(&data); err == nil {
							result := songModel.Update(map[string]string{
								"id":           id,
								"title":        song,
								"group_id":     groupId,
								"link":         data["link"],
								"release_date": data["release_date"],
								"text":         data["text"],
							}, id)
							if _, ok := result["id"]; !ok {
								customLog.Logging(errors.New(utils.ConcatSlice([]string{
									"error: error when trying to enrich a song",
									"ID: ",
									id,
								})))
							}
						} else {
							customLog.Logging(err)
						}
					}
				}
			}
		} else {
			response["error"] = "Check parameters"
		}
	}
	json.NewEncoder(w).Encode(response)
}

func (router *Router) deleteSong(w http.ResponseWriter, r *http.Request) {
	params := router.initProcess(w, r, false)
	id, err := strconv.Atoi(params["id"])
	fmt.Println(id)
	if err != nil {
		customLog.Logging(err)
	} else {
		songModel := (*&models.Song{}).Init()
		response = map[string]interface{}{"data": songModel.Delete(id)}
	}
	json.NewEncoder(w).Encode(response)
}
