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
	//list, err := db.GetList("ttt", []string{}, map[string]interface{}{"Id": 4, "Name": "helloworld"})
	//if err != nil {
	//	log.Fatal(err.Error())
	//}
	//fmt.Println(list)
	rowsAffected, err := db.Delete("ttt", map[string]interface{}{})
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(rowsAffected)

}
