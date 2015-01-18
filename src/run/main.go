package main

import (
	"controllers"
	"fmt"
	"libraries/common"
	"log"
	"os"
	"reflect"
	//"runtime"
	"strings"
)

func main() {
	/*
	 * 使用方法run [class/method] param1 param2 param3 ...
	 * class 指定的是controllers包下面的一个类
	 * method 指定的是controllers包下面的类中的方式
	 * param 列表是controllers包下面的类方法所需要的参数列表。
	 * 注意：由于从控制台传入的参数列表都是string类型，因此controller类方法中接收的参数也需要是string类型
	 */

	var class string
	var method string
	if len(os.Args) < 2 {
		fmt.Println("Usage run [Class/Method] param1 param2 param3 ...")
		os.Exit(1)
	}

	//移除文件名
	input := os.Args[1:]

	//第一个参数即是逻辑定位参数
	logic := input[0]
	splited := strings.Split(logic, "/")
	if len(splited) < 2 {
		class = strings.ToLower(splited[0])
		method = "Index"
	} else {
		class = strings.ToLower(splited[0])
		method = splited[1]
	}

	//设置可用的cpu核心数
	common.SetCPUNum()

	defer func() {
		if r := recover(); r != nil {
			/*var buf []byte = make([]byte, 1024)
			c := runtime.Stack(buf, false)
			fmt.Println(string(buf[0:c]))
			*/
			log.Fatalf("Runtime error caught: %v", r)
		}
	}()

	ctrls := make(map[string]interface{})

	//setting controller structs
	ctrls["testing"] = &controllers.Testing{}
	ctrls["webserver"] = &controllers.Webserver{}

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

	//后面的参数是逻辑所需要的参数，将依次传入到方法中
	input = input[1:]

	var in []reflect.Value = make([]reflect.Value, len(input))
	for i, v := range input {
		in[i] = reflect.ValueOf(v)
	}
	f.Call(in)

}
