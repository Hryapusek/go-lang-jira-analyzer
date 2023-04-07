package main

import (
	"ApiServer/internals/config"
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"syscall"
	"time"
)

type route struct {
	prefix  string
	target  *url.URL
	timeout time.Duration // время ожидания
}

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

	// Разрешенные маршруты
	routes := []route{
		{cfg.AnalyticsAPIPrefix, analyticsTarget, time.Duration(cfg.AnalyticsTimeout) * time.Second},
		{cfg.ResourceAPIPrefix, resourceTarget, time.Duration(cfg.ResourceTimeout) * time.Second},
		{cfg.ConnectorAPIPrefix, connectorTarget, 0}, // микросервис connector без лимита
	}

	mux := http.NewServeMux()
	mux.HandleFunc(cfg.GatewayAPIPrefix, func(w http.ResponseWriter, r *http.Request) {
		for _, rt := range routes {
			if strings.HasPrefix(r.URL.Path, rt.prefix) {
				proxy := httputil.NewSingleHostReverseProxy(rt.target)

				if rt.timeout > 0 {
					ctx, cancel := context.WithTimeout(r.Context(), rt.timeout)
					defer cancel()
					r = r.WithContext(ctx)

					done := make(chan struct{})
					go func() {
						proxy.ServeHTTP(w, r)
						close(done)
					}()

					select {
					case <-done:
					case <-ctx.Done():
						w.WriteHeader(http.StatusGatewayTimeout)
						_, err := fmt.Fprint(w, "Gateway Timeout")
						if err != nil {
							return
						}
					}
				} else {
					proxy.ServeHTTP(w, r)
				}
				return
			}
		}
		http.Error(w, "Not found", http.StatusNotFound)
	})

	stop := make(chan os.Signal, 1)

	fmt.Println("Starting reverse proxy server on ", gatewayAddress)
	go func() {
		if err := http.ListenAndServe(gatewayAddress, mux); err != nil {
			fmt.Printf("ListenAndServe error: %v\n", err)
			stop <- syscall.SIGTERM
		}
	}()
}
