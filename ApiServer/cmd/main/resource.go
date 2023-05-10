package main

import (
	"ApiServer/internals/config"
	endpoints "ApiServer/internals/endpoints/Resource"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	file, err := os.Create("./log/resourceLog.txt")
	if err != nil {
		log.SetOutput(os.Stdout)
		log.Println("Cannot create log file", err)
	} else {
		logsOutput := io.MultiWriter(os.Stdout, file)
		log.SetOutput(logsOutput)

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Println("Unable to close log file")
			}
		}(file)
	}

	cfg := config.LoadResourceConfig("configs/server.yaml")
	log.Printf("Create handler for mask \"%s\"", cfg.MainAPIPrefix+cfg.ResourceAPIPrefix)

	resourceRouter := mux.NewRouter()
	resourceRouter.HandleFunc(cfg.MainAPIPrefix+cfg.ResourceAPIPrefix+"issues/{id:[0-9]+}", endpoints.HandlerGetIssue).Methods("GET")
	resourceRouter.HandleFunc(cfg.MainAPIPrefix+cfg.ResourceAPIPrefix+"histories/{id:[0-9]+}", endpoints.HandlerGetHistory).Methods("GET")
	resourceRouter.HandleFunc(cfg.MainAPIPrefix+cfg.ResourceAPIPrefix+"projects/{id:[0-9]+}", endpoints.HandlerGetProject).Methods("GET")

	resourceRouter.HandleFunc(cfg.MainAPIPrefix+cfg.ResourceAPIPrefix+"issues/", endpoints.HandlerGetIssuesByProjectId).
		Methods("GET").
		Queries("projectId", "{projectId:[0-9]+}").
		Queries("offset", "{offset?:[0-9]+}").
		Queries("limit", "{limit?:[0-9]+}")

	resourceRouter.HandleFunc(cfg.MainAPIPrefix+cfg.ResourceAPIPrefix+"projects/", endpoints.HandlerGetAllProject).
		Methods("GET").
		Queries("offset", "{offset?:[0-9]+}").
		Queries("limit", "{limit?:[0-9]+}")

	resourceRouter.HandleFunc(cfg.MainAPIPrefix+cfg.ResourceAPIPrefix+"get_project_by_title",
		endpoints.HandlerGetProjectByTitle).Methods("GET").Queries("title", "{title}")

	resourceRouter.HandleFunc(cfg.MainAPIPrefix+cfg.ResourceAPIPrefix+"issues/", endpoints.HandlerPostIssue).Methods("POST")
	resourceRouter.HandleFunc(cfg.MainAPIPrefix+cfg.ResourceAPIPrefix+"histories/", endpoints.HandlerPostHistory).Methods("POST")
	resourceRouter.HandleFunc(cfg.MainAPIPrefix+cfg.ResourceAPIPrefix+"projects/", endpoints.HandlerPostProject).Methods("POST")

	resourceAddress := fmt.Sprintf("%s:%d", cfg.ResourceHost, cfg.ResourcePort)

	log.Printf("Start resource server at %s", resourceAddress)
	err = http.ListenAndServe(resourceAddress, resourceRouter)
	if err != nil {
		log.Fatalf("Unable to start resource server at %s, because of %s", resourceAddress, err.Error())
	}
}
