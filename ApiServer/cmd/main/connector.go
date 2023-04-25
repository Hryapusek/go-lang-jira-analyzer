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
	file, err := os.Create("./log/connectorLog.txt")
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

	cfg := config.LoadConnectorConfig("configs/server.yaml")
	log.Printf("Create handler for mask \"%s\"", cfg.MainAPIPrefix+cfg.ConnectorAPIPrefix)

	connectorRouter := http.NewServeMux()
	connectorRouter.HandleFunc(cfg.MainAPIPrefix+cfg.ConnectorAPIPrefix, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("New request for connector server at %s", r.URL.Path)
		w.WriteHeader(200)
		_, err := w.Write([]byte("Status is OK"))
		if err != nil {
			return
		}
	})

	connectorAddress := fmt.Sprintf("%s:%d", cfg.ConnectorHost, cfg.ConnectorPort)

	log.Printf("Start connector server at %s", connectorAddress)
	err = http.ListenAndServe(connectorAddress, connectorRouter)
	if err != nil {
		log.Fatalf("Unable to start connector server at %s", connectorAddress)
	}
}
