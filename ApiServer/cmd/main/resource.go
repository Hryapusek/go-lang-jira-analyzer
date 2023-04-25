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

	resourceRouter := http.NewServeMux()
	resourceRouter.HandleFunc(cfg.MainAPIPrefix+cfg.ResourceAPIPrefix, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("New request for resource server at %s", r.URL.Path)
		w.WriteHeader(200)
		_, err := w.Write([]byte("Status is OK"))
		if err != nil {
			return
		}
	})

	resourceAddress := fmt.Sprintf("%s:%d", cfg.ResourceHost, cfg.ResourcePort)

	log.Printf("Start resource server at %s", resourceAddress)
	err = http.ListenAndServe(resourceAddress, resourceRouter)
	if err != nil {
		log.Fatalf("Unable to start resource server at %s", resourceAddress)
	}
}
