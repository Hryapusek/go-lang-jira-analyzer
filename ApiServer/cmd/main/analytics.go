package main

import (
	"ApiServer/internals/config"
	"ApiServer/internals/endpoints/Analytics"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	file, err := os.Create("./log/analyticsLog.txt")
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

	cfg := config.LoadAnalyticsConfig("configs/server.yaml")

	analyticsRouter := mux.NewRouter()

	log.Printf("Create handler for mask \"%s\"", cfg.MainAPIPrefix+cfg.AnalyticsAPIPrefix+"services")
	analyticsRouter.HandleFunc(cfg.MainAPIPrefix+cfg.AnalyticsAPIPrefix+"services", endpoints.AnalyticsServices)
	log.Printf("Create handler for mask \"%s\"", cfg.MainAPIPrefix+cfg.AnalyticsAPIPrefix+"{group:[1-5]}")
	analyticsRouter.HandleFunc(cfg.MainAPIPrefix+cfg.AnalyticsAPIPrefix+"{group:[1-5]}",
		endpoints.GetGraph).Queries("project", "{projectName}").Methods("GET") // если не указывать метод Queries, обрабтываются заранее невалидные запросы по аргументам
	// т.е. -- оставить напоминание.

	log.Printf("Create handler for mask \"%s\"", cfg.MainAPIPrefix+cfg.AnalyticsAPIPrefix)
	analyticsRouter.HandleFunc(cfg.MainAPIPrefix+cfg.AnalyticsAPIPrefix, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("New request for analytics server at %s", r.URL.Path)
		w.WriteHeader(404)
		_, err := w.Write([]byte("Status is 404"))
		if err != nil {
			return
		}
	})

	analyticsAddress := fmt.Sprintf("%s:%d", cfg.AnalyticsHost, cfg.AnalyticsPort)

	log.Printf("Start analytics server at %s", analyticsAddress)
	err = http.ListenAndServe(analyticsAddress, analyticsRouter)
	if err != nil {
		log.Fatalf("Unable to start connector server at %s", analyticsAddress)
	}
}
