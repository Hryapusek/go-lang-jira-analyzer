package endpoints

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var services = []string{
	"/api/v1/graph/services",
	"/api/v1/graph/1",
	"/api/v1/graph/2",
	"/api/v1/graph/3",
	"/api/v1/graph/4",
	"/api/v1/graph/5",
	"/api/v1/graph/6",
}

func AnalyticsServices(rw http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(services)
	if err != nil {
		log.Fatalf("Error with JSON response on request %s", r.URL.Path)
	}

	_, err = rw.Write(data)
	log.Printf("Writed data on \"/api/v1/graph/services\" request")
	if err != nil {
		return
	}
}

func GetGraph(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	group, err := strconv.Atoi(vars["group"])
	if err != nil {
		log.Printf("invalid group request in path \"%s\"", r.URL.Path)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	projectId, err := strconv.Atoi(r.URL.Query().Get("project"))
	if err != nil {
		log.Printf("invalid group request in path \"%s\"", r.URL.Path)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Incoming request on endpoint /api/v1/graph/%d?project=%s", group, projectId)

	var data []byte

	switch group {
	case 1:
		data, err = json.MarshalIndent(GraphOne(projectId), "", "\t")
	case 2:
		data, err = json.MarshalIndent(GraphTwo(projectId), "", "\t")
	case 3:
		data, err = json.MarshalIndent(GraphThree(projectId), "", "\t")
	case 4:
		data, err = json.MarshalIndent(GraphFour(projectId), "", "\t")
	case 5:
		data, err = json.MarshalIndent(GraphFive(projectId), "", "\t")
	case 6:
		data, err = json.MarshalIndent(GraphSix(projectId), "", "\t")
	default:
		log.Printf("Not exisiting group parameter at /api/v1/graph/{group:[1-6]}. Ты как сюда попал??")
		rw.WriteHeader(403)
		return
	}

	if err != nil {
		log.Printf("Internal error with marshaling data from /api/v1/graph/%d?project=%s", group, projectId)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = rw.Write(data)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}
