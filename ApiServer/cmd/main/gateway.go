package main

import (
	"ApiServer/internals/config"
	"context"
	_ "context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	_ "strings"
	_ "syscall"
	"time"
)

func main() {
	cfg := config.LoadGatewayConfig("configs/server.yaml")

	analyticsTarget, err := url.Parse(fmt.Sprintf("http://%s:%d", cfg.AnalyticsHost, cfg.AnalyticsPort))
	if err != nil {
		fmt.Printf("Error parsing target URL: %v\n", err)
		os.Exit(1)
	}

	resourceTarget, err := url.Parse(fmt.Sprintf("http://%s:%d", cfg.ResourceHost, cfg.ResourcePort))
	if err != nil {
		fmt.Printf("Error parsing target URL: %v\n", err)
		os.Exit(1)
	}

	connectorTarget, err := url.Parse(fmt.Sprintf("http://%s:%d", cfg.ConnectorHost, cfg.ConnectorPort))
	if err != nil {
		fmt.Printf("Error parsing target URL: %v\n", err)
		os.Exit(1)
	}

	gatewayAddress := fmt.Sprintf("%s:%d", cfg.GatewayHost, cfg.GatewayPort)

	gatewayMux := http.NewServeMux()

	gatewayMux.HandleFunc(cfg.GatewayAPIPrefix+cfg.AnalyticsAPIPrefix, func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cfg.AnalyticsTimeout))
		defer cancel()

		r = r.WithContext(ctx)

		proxy := httputil.NewSingleHostReverseProxy(analyticsTarget)
		proxy.ServeHTTP(w, r)
	})

	gatewayMux.HandleFunc(cfg.GatewayAPIPrefix+cfg.ResourceAPIPrefix, func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cfg.ResourceTimeout))
		defer cancel() // нужно добавить логирование запросов + вывод истекших по timeout

		r = r.WithContext(ctx)

		proxy := httputil.NewSingleHostReverseProxy(resourceTarget)
		proxy.ServeHTTP(w, r)
	})

	gatewayMux.HandleFunc(cfg.GatewayAPIPrefix+cfg.ConnectorAPIPrefix, func(w http.ResponseWriter, r *http.Request) {
		proxy := httputil.NewSingleHostReverseProxy(connectorTarget)
		proxy.ServeHTTP(w, r)
	})

	err = http.ListenAndServe(gatewayAddress, gatewayMux)
	if err != nil {
		return
	}
}
