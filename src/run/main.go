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
	if len(os.Args) != 3 {
		fmt.Println("Usage run [class] [Method]")
		os.Exit(1)
	}

	class := os.Args[1]
	method := os.Args[2]

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

	action, ok := ctrls[class]

	if !ok {
		fmt.Println("Class(", class, ") is invalid.")
		os.Exit(1)
	}

	f := reflect.ValueOf(action).MethodByName(method)
	if !f.IsValid() {
		fmt.Println("Method '", method, "' is invalid.")
		os.Exit(1)
	}
	var in []reflect.Value = make([]reflect.Value, 0)
	f.Call(in)

}
