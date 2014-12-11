package main

import (
	"fmt"
	"reflect"
)

type Bird struct {
	Name           string
	LifeExpectance int
}

func main() {
	var sparrow *Bird = &Bird{"Sparrow", 3}
	v := reflect.ValueOf(sparrow).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fmt.Printf("var %s %s = %v \n", t.Field(i).Name, v.Field(i).Type(), v.Field(i).Interface())
	}
}

/*
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
*/
