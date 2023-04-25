package main

import (
	"ApiServer/internals/config"
	"context"
	_ "context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	_ "strings"
	_ "syscall"
	"time"
)

func main() {
	file, err := os.Create("./log/gatewayLog.txt")
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

	cfg := config.LoadGatewayConfig("configs/server.yaml")

	analyticsTarget, err := url.Parse(fmt.Sprintf("http://%s:%d", cfg.AnalyticsHost, cfg.AnalyticsPort))
	if err != nil {
		log.Fatalf("Error parsing target URL: %v\n", err)
	}

	resourceTarget, err := url.Parse(fmt.Sprintf("http://%s:%d", cfg.ResourceHost, cfg.ResourcePort))
	if err != nil {
		log.Fatalf("Error parsing target URL: %v\n", err)
	}

	connectorTarget, err := url.Parse(fmt.Sprintf("http://%s:%d", cfg.ConnectorHost, cfg.ConnectorPort))
	if err != nil {
		log.Fatalf("Error parsing target URL: %v\n", err)
	}

	gatewayAddress := fmt.Sprintf("%s:%d", cfg.GatewayHost, cfg.GatewayPort)

	gatewayMux := http.NewServeMux()

	gatewayMux.HandleFunc(cfg.GatewayAPIPrefix+cfg.AnalyticsAPIPrefix, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Recieved request for analytics proxy: %s", r.URL.Path)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cfg.AnalyticsTimeout))
		defer cancel()

		r = r.WithContext(ctx)

		proxy := httputil.NewSingleHostReverseProxy(analyticsTarget)
		proxy.ServeHTTP(w, r)
	})

	gatewayMux.HandleFunc(cfg.GatewayAPIPrefix+cfg.ResourceAPIPrefix, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Recieved request for resource proxy: %s", r.URL.Path)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cfg.ResourceTimeout))
		defer cancel() // нужно добавить логирование запросов + вывод истекших по timeout

		r = r.WithContext(ctx)

		proxy := httputil.NewSingleHostReverseProxy(resourceTarget)
		proxy.ServeHTTP(w, r)
	})

	gatewayMux.HandleFunc(cfg.GatewayAPIPrefix+cfg.ConnectorAPIPrefix, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Recieved request for connector proxy: %s", r.URL.Path)
		proxy := httputil.NewSingleHostReverseProxy(connectorTarget)
		proxy.ServeHTTP(w, r)
	})

	err = http.ListenAndServe(gatewayAddress, gatewayMux)
	if err != nil {
		log.Fatalf("Unable to start gateway server at %s", gatewayAddress)
	}
}
