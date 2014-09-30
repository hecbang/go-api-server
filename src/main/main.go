package main

import (
	"controllers/openapi"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		var api = &openapi.OpenApi{}
		api.SetRequestData(w, r)
		api.Index(w, r)
	})
	log.Fatal(http.ListenAndServe(":9090", nil))
}
