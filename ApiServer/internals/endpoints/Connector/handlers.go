package Connector

import (
	"fmt"
	"net/http"
)

func All(rw http.ResponseWriter, r *http.Response) {
	rw.WriteHeader(200)
	_, err := rw.Write([]byte(fmt.Sprintf("Here should be graph endpoint for groups")))
	if err != nil {
		return
	}
}
