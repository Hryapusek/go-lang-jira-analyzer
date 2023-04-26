package endpoints

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
)

func GetIssue(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(200)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid Issue ID in path \"%s\"", r.URL.Path)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(GetIssueInfoByID(id))
	if err != nil {
		log.Printf("Error with extracting info about issue project with id=%d", id)
		rw.WriteHeader(400)
		return
	}

	rw.WriteHeader(200)
	_, err = rw.Write(data)
	if err != nil {
		return
	}
}

func GetHistory(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(200)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid Issue ID in path \"%s\"", r.URL.Path)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(GetHistoryInfoByID(id))
	if err != nil {
		log.Printf("Error with extracting info about history with id=%d", id)
		rw.WriteHeader(400)
		return
	}

	rw.WriteHeader(200)
	_, err = rw.Write(data)
	if err != nil {
		return
	}
}

func GetProject(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(200)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid Issue ID in path \"%s\"", r.URL.Path)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(GetProjectInfoByID(id))
	if err != nil {
		log.Printf("Error with extracting info about project with id=%d", id)
		rw.WriteHeader(400)
		return
	}

	rw.WriteHeader(200)
	_, err = rw.Write(data)
	if err != nil {
		return
	}

}

func PostIssue(rw http.ResponseWriter, r *http.Request) {
	var data IssueInfo
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// Обработка ошибки
	}
	if err := json.Unmarshal(body, &data); err != nil {
		// Обработка ошибки
	}

	err = PutIssueToDB(data)
	if err != nil {
		log.Printf("Error %s occured while puting issue to DB", err.Error())
	} else {
		rw.WriteHeader(http.StatusCreated)
	}
}

func PostHistory(rw http.ResponseWriter, r *http.Request) {
	var data HistoryInfo
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// Обработка ошибки
	}
	if err := json.Unmarshal(body, &data); err != nil {
		// Обработка ошибки
	}

	err = PutHistoryToDB(data)
	if err != nil {
		log.Printf("Error %s occured while puting history to DB", err.Error())
	} else {
		rw.WriteHeader(http.StatusCreated)
	}
}

func PostProject(rw http.ResponseWriter, r *http.Request) {
	var data ProjectInfo
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// Обработка ошибки
	}
	if err := json.Unmarshal(body, &data); err != nil {
		// Обработка ошибки
	}

	err = PutProjectToDB(data)
	if err != nil {
		log.Printf("Error %s occured while puting project to DB", err.Error())
	} else {
		rw.WriteHeader(http.StatusCreated)
	}
}
