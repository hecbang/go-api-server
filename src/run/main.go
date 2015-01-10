package main

import (
	"fmt"
	"libraries/common"
	"log"
	"runtime"
)

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
	fmt.Println("get here")
}
