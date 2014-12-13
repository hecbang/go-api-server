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
	filepath := path.CONFIG_PATH + "db.json"
	content, err := common.Read(filepath)
	if err != nil {
		log.Fatal(err.Error())
	}
	var v map[string]DbItem
	if err := json.Unmarshal(content, &v); err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(v["default"].Slave[common.Rand(0, len(v["default"].Slave)-1)])
}
