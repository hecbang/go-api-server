package main

import (
	"const/path"
	"encoding/json"
	"fmt"
	"libraries/common"
	"log"
	"runtime"
)

type DbItem struct {
	Master string
	Slave  []string
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			var buf []byte = make([]byte, 1024)
			c := runtime.Stack(buf, false)
			fmt.Println(string(buf[0:c]))
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

	db := common.NewMySql()
	//list, err := db.GetList("teacher", []string{}, map[string]string{"Id": "4;2;3"})
	//if err != nil {
	//	log.Fatal(err.Error())
	//}
	//fmt.Println(list)
	id, err := db.Insert("ttt", map[string]interface{}{"Number": 12, "Name": "yorkershi"})
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(id)

}
