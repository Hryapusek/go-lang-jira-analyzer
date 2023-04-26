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
	resourceRouter.HandleFunc(cfg.MainAPIPrefix+cfg.ResourceAPIPrefix+"/issue/{id:[0-9]+}", endpoints.GetIssue).Methods("GET")
	resourceRouter.HandleFunc(cfg.MainAPIPrefix+cfg.ResourceAPIPrefix+"/history/{id:[0-9]+}", endpoints.GetHistory).Methods("GET")
	resourceRouter.HandleFunc(cfg.MainAPIPrefix+cfg.ResourceAPIPrefix+"/project/{id:[0-9]+}", endpoints.GetProject).Methods("GET")

	resourceRouter.HandleFunc(cfg.MainAPIPrefix+cfg.ResourceAPIPrefix+"/issue/{id:[0-9]+}", endpoints.PostIssue).Methods("POST")
	resourceRouter.HandleFunc(cfg.MainAPIPrefix+cfg.ResourceAPIPrefix+"/history/{id:[0-9]+}", endpoints.PostHistory).Methods("POST")
	resourceRouter.HandleFunc(cfg.MainAPIPrefix+cfg.ResourceAPIPrefix+"/project/{id:[0-9]+}", endpoints.PostProject).Methods("POST")

	resourceAddress := fmt.Sprintf("%s:%d", cfg.ResourceHost, cfg.ResourcePort)

	log.Printf("Start resource server at %s", resourceAddress)
	err = http.ListenAndServe(resourceAddress, resourceRouter)
	if err != nil {
		log.Fatalf("Unable to start resource server at %s, because of %s", resourceAddress, err.Error())
	}
}
