package main

import (
	"const/path"
	"fmt"
	"libraries/common"
	"os"
	"log"
)

func main() {
	filepath := path.CONFIG_PATH + "db.json"
	fp, err := os.Open(filepath)
	if err != nil {
		
	}
	fmt.Println(common.FileExists(filepath))
}
