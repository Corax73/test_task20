package extapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"songLibrary/customLog"
	"songLibrary/utils"
	"strings"
	"time"

	randomDataTime "github.com/duktig-solutions/go-random-date-generator"
	"github.com/gorilla/mux"
)

type ExtRouter struct {
	*mux.Router
}

func (router *ExtRouter) Init() *ExtRouter {
	r := mux.NewRouter()
	r.HandleFunc("/info/{group:[a-zA-Z0-9\\s]+}/{song:[a-zA-Z0-9\\s]+}", router.getOneSongs).Methods("GET")
	return &ExtRouter{r}
}

func (router *ExtRouter) initProcess(w http.ResponseWriter, r *http.Request, isPost bool) map[string]string {
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

func (router *ExtRouter) checkEnv() bool {
	var resp bool
	envData := utils.GetConfFromEnvFile()
	if val, ok := envData["CONSOLE_OUT"]; ok && val == "true" {
		resp = true
	}
	return resp
}

func (router *ExtRouter) consoleOutput(r *http.Request) {
	fmt.Println(strings.Join([]string{time.Now().Format(time.RFC3339), r.Method, r.RequestURI, r.UserAgent()}, " "))
}

func (router *ExtRouter) getOneSongs(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]string)
	params := router.initProcess(w, r, false)
	if group, ok := params["group"]; ok {
		if song, ok := params["song"]; ok {
			response["release_date"] = "2024-11-28"
			randomDate, err := randomDataTime.GenerateDate("2000-08-01", "2024-08-01")
			if err != nil {
				customLog.Logging(err)
			} else {
				response["release_date"] = randomDate
			}
			response["text"] = "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"
			response["link"] = utils.ConcatSlice([]string{"/", group, "/", song})
		}
	}
	json.NewEncoder(w).Encode(response)
}
