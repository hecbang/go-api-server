package main

import (
	"const/path"
	"fmt"
	"libraries/common"
	"os"
)

func main() {
	filepath := path.CONFIG_PATH + "db.json"
	fp, err := os.Open(filepath)
	if err != nil {
		log.
	}
	fmt.Println(common.FileExists(filepath))
}
