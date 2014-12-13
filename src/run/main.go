package main

import (
	"const/path"
	"encoding/json"
	"fmt"
	"libraries/common"
	"log"
)

type DbItem struct {
	Master string
	Slave  []string
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Runtime error caught: %v", r)
		}
	}()

	filepath := path.CONFIG_PATH + "db.json"
	content, err := common.Read(filepath)
	if err != nil {
		log.Fatal(err.Error())
	}
	var v map[string]DbItem
	if err := json.Unmarshal(content, &v); err != nil {
		log.Fatal(err.Error())
	}

	db := common.NewQuery()
	fmt.Println(db)
}
