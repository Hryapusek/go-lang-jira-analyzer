package main

import (
	"ApiServer/internals/config"
	"fmt"
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

	log.Printf("Create handler for mask \"%s\"", cfg.MainAPIPrefix+cfg.AnalyticsAPIPrefix)

	analyticsRouter := http.NewServeMux()
	analyticsRouter.HandleFunc(cfg.MainAPIPrefix+cfg.AnalyticsAPIPrefix, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("New request for connector server at %s", r.URL.Path)
		w.WriteHeader(200)
		_, err := w.Write([]byte("Status is OK"))
		if err != nil {
			return
		}
	})

	analyticsAddress := fmt.Sprintf("%s:%d", cfg.AnalyticsHost, cfg.AnalyticsPort)

	log.Printf("Start connector server at %s", analyticsAddress)
	err = http.ListenAndServe(analyticsAddress, analyticsRouter)
	if err != nil {
		log.Fatalf("Unable to start connector server at %s", analyticsAddress)
	}
}
