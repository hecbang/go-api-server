package main

import (
	"const/path"
	"fmt"
	"libraries/common"
)

func main() {
	filepath := path.CONFIG_PATH + "db.json"

	fmt.Println(common.FileExists(filepath))
}
