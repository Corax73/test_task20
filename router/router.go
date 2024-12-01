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
	return &Router{r}
}

func (router *Router) initProcess(w http.ResponseWriter, r *http.Request) map[string]string {
	w.Header().Set("Content-Type", "application/json")
	if router.checkEnv() {
		router.consoleOutput(r)
	}
	return mux.Vars(r)
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
	params := router.initProcess(w, r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		customLog.Logging(err)
	} else {
		songModel := (*&models.Song{}).Init()
		response = map[string]interface{}{"data": songModel.GetOne(id)}
	}
	json.NewEncoder(w).Encode(response)
}
