package endpoints

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func graph_placeholder(projectName string) map[int]int { // для проверки функционала,
	// TODO: снести в конце
	return map[int]int{
		1: 1,
		2: 4,
	}
}

var services = []string{
	"/api/v1/graph/services",
	"/api/v1/graph/1",
	"/api/v1/graph/2",
	"/api/v1/graph/3",
	"/api/v1/graph/4",
	"/api/v1/graph/5",
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
	projectName := r.URL.Query().Get("project")

	log.Printf("Incoming request on endpoint /api/v1/graph/%d?project=%s", group, projectName)

	var data []byte

	switch group {
	case 1:
		data, err = json.Marshal(GraphOne(projectName))
	case 2:
		data, err = json.Marshal(GraphTwo(projectName))
	case 3:
		data, err = json.Marshal(GraphThree(projectName))
	case 4:
		data, err = json.Marshal(graph_placeholder(projectName))
	case 5:
		data, err = json.Marshal(graph_placeholder(projectName))
	}

	if err != nil {
		log.Printf("Internal error with marshaling data from /api/v1/graph/%d?project=%s", group, projectName)
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
