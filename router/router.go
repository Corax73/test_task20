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

// @var response stub for answers on routes.
var response map[string]interface{}

func (router *Router) Init() *Router {
	r := mux.NewRouter()
	r.HandleFunc("/songs/", router.getSongs).Methods("GET")
	r.HandleFunc("/songs/{id:[0-9]+}", router.getOneSongs).Methods("GET")
	r.HandleFunc("/songs/{id:[0-9]+}/couplet/{couplet_number:[0-9]+}", router.getOneCouplet).Methods("GET")
	r.HandleFunc("/songs/", router.createSong).Methods("POST")
	r.HandleFunc("/songs/{id:[0-9]+}", router.updateSong).Methods("PUT")
	r.HandleFunc("/songs/{id:[0-9]+}", router.deleteSong).Methods("DELETE")
	r.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir("./swagger/"))))
	return &Router{r}
}

// initProcess returns a map of request parameters, causes console output on request.
func (router *Router) initProcess(w http.ResponseWriter, r *http.Request, getPost bool) map[string]string {
	var resp map[string]string
	w.Header().Set("Content-Type", "application/json")
	if router.checkEnv() {
		router.consoleOutput(r)
	}
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
	if getPost {
		err := json.NewDecoder(r.Body).Decode(&resp)
		if err != nil {
			customLog.Logging(err)
		}
	}
	return resp
}

// checkEnv looks for a key `CONSOLE_OUT` in the .env file and returns true if its value is true.
func (router *Router) checkEnv() bool {
	var resp bool
	envData := utils.GetConfFromEnvFile()
	if val, ok := envData["CONSOLE_OUT"]; ok && val == "true" {
		resp = true
	}
	return resp
}

// consoleOutput displays the time, route and request method to the console.
func (router *Router) consoleOutput(r *http.Request) {
	fmt.Println(strings.Join([]string{time.Now().Format(time.RFC3339), r.Method, r.RequestURI, r.UserAgent()}, " "))
}

// getSongs returns a list of entities, can use limit and offset parameters.
func (router *Router) getSongs(w http.ResponseWriter, r *http.Request) {
	params := router.initProcess(w, r, false)
	songModel := (*&models.Song{}).Init()
	response = map[string]interface{}{"data": songModel.GetList(params)}
	json.NewEncoder(w).Encode(response)
}

// getOneSongs returns entity data by parameter `id`.
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

// createSong by post parameters creates an entity, you, if the model implements the interface,
// then a request is made to enrich the entity data.
func (router *Router) createSong(w http.ResponseWriter, r *http.Request) {
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

// deleteSong deletes an entity using the parameter `id`.
func (router *Router) deleteSong(w http.ResponseWriter, r *http.Request) {
	params := router.initProcess(w, r, false)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		customLog.Logging(err)
	} else {
		songModel := (*&models.Song{}).Init()
		response = map[string]interface{}{"data": songModel.Delete(id)}
	}
	json.NewEncoder(w).Encode(response)
}

// getOneCouplet returns the text of the song by verses using parameters `id` and `couplet_number`.
func (router *Router) getOneCouplet(w http.ResponseWriter, r *http.Request) {
	params := router.initProcess(w, r, false)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		customLog.Logging(err)
	} else {
		couplet_number, err := strconv.Atoi(params["couplet_number"])
		if err != nil {
			customLog.Logging(err)
		} else {
			songModel := (*&models.Song{}).Init()
			response = map[string]interface{}{"data": songModel.GetOneCouplet(id, couplet_number)}
		}
	}
	json.NewEncoder(w).Encode(response)
}

// updateSong updates entity.
func (router *Router) updateSong(w http.ResponseWriter, r *http.Request) {
	params := router.initProcess(w, r, true)
	if _, ok := params["id"]; ok {
		songModel := (*&models.Song{}).Init()
		response = map[string]interface{}{"data": songModel.Update(params, params["id"])}
	}
	json.NewEncoder(w).Encode(response)
}
