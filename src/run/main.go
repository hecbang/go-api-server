package main

import (
	"controllers"
	"fmt"
	"libraries/common"
	"log"
	"os"
	"reflect"
	"runtime"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage run [Method]")
		os.Exit(1)
	}

	common.SetCPUNum()
	defer func() {
		if r := recover(); r != nil {
			var buf []byte = make([]byte, 1024)
			c := runtime.Stack(buf, false)
			fmt.Println(string(buf[0:c]))
			log.Fatalf("Runtime error caught: %v", r)
		}
	}()

	ctrls := make(map[string]interface{})

	//setting controller structs
	ctrls["testing"] = &controllers.Testing{}

	m := os.Args[1]
	f := reflect.ValueOf(ctrls["testing"]).MethodByName(m)
	if !f.IsValid() {
		fmt.Println("Method '", m, "' is invalid.")
		os.Exit(1)
	}
	var in []reflect.Value = make([]reflect.Value, 0)
	f.Call(in)

}
