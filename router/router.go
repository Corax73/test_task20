package router

import (
	"encoding/json"
	"net/http"
	"songLibrary/customLog"
	"songLibrary/models"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Router struct {
	*mux.Router
}

func (router *Router) Init() *Router {
	r := mux.NewRouter()
	r.HandleFunc("/songs/{id}", router.getOneSongs).Methods("GET")
	return &Router{r}
}

func (router *Router) getOneSongs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	resp := map[string]interface{}{"id": "123",
		"group_id":    "321",
		"title":       "test",
		"releaseDate": time.Now().String(),
		"text":        "text`",
		"link":        "/",
	}
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		customLog.Logging(err)
	} else {
		songModel := (*&models.Song{}).Init()
		resp = songModel.GetOne(id)
	}
	json.NewEncoder(w).Encode(resp)
}
