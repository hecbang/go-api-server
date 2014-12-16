package main

import (
	"fmt"
	"libraries/common"
	"log"
	"runtime"
	"time"
)

type DbItem struct {
	Master string
	Slave  []string
}

func main() {
	common.SetCPUNum()
	defer func() {
		if r := recover(); r != nil {
			var buf []byte = make([]byte, 1024)
			c := runtime.Stack(buf, false)
			fmt.Println(string(buf[0:c]))
			log.Fatalf("Runtime error caught: %v", r)

		}
	}()
	start := time.Now().UnixNano()

	db := common.NewMySql()
	//db.Begin()
	//list, err := db.GetList("ttt", []string{}, map[string]interface{}{"Id": 4, "Name": "helloworld"})
	//if err != nil {
	//	log.Fatal(err.Error())
	//}
	//fmt.Println(list)
	lastInsertId, err := db.Insert("ttt", map[string]interface{}{"Number": 43, "Name": "yorkershi", "Point": 8.8})
	if err != nil {
		log.Fatal(err.Error())
	}
	info, err1 := db.GetDictionary("ttt", []string{}, map[string]interface{}{"Id": lastInsertId})
	if err1 != nil {
		log.Fatalln(err1.Error())
	}
	fmt.Println(info)
	//db.Commit()
	fmt.Println(lastInsertId)
	end := time.Now().UnixNano()
	fmt.Println(common.Round((float64(end)-float64(start))/1000/1000, 3), "ms")

}
