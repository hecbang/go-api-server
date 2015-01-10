package main

import (
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
	common.SetCPUNum()
	defer func() {
		if r := recover(); r != nil {
			var buf []byte = make([]byte, 1024)
			c := runtime.Stack(buf, false)
			fmt.Println(string(buf[0:c]))
			log.Fatalf("Runtime error caught: %v", r)

		}
	}()

	filename := "CommonAutoSpeechNotice.ddp"
	slice, err := common.ReadLine(filename)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(slice[50])

	var jsonstring string = `{"Action":"UserManagement", "Method":"auth", "SystemId":"15"}`
	v := make(map[string]interface{})
	//var vv struct {
	//	Action   string
	//	Method   string
	//	SystemId string
	//}

	err1 := json.Unmarshal([]byte(jsonstring), &v)
	if err1 != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(v)
}
