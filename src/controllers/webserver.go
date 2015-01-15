package controllers

import (
	"fmt"
	"log"
	"net/http"
)

type Webserver struct {
}

func (this *Webserver) Index() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "helloyorkershi")
	})
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
