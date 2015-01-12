package controllers

import (
	"const/path"
	"fmt"
	"libraries/common"
	"log"
	"models/testing"
)

type Testing struct {
}

func (this *Testing) Concurrence() {
	jq := common.LoadJsonQuery(path.CONFIG_PATH + "testing" + path.DS + "concurrence.json")

	db := common.NewMySqlInstance("testdata")

	name, err := jq.String("group", "name")
	if err != nil {
		log.Fatalln(err.Error())
	}
	parameter, err := jq.String("group", "parameter")
	if err != nil {
		log.Fatalln(err.Error())
	}
	data := map[string]interface{}{
		"Name":              name,
		"SettingParameters": parameter,
		"LogTime":           common.Date("Y-m-d H:i:s"),
	}

	lastid, err := db.Insert("db_group", data)

	c := start
	for c < max {
		testing.Concurrence(config["concurrence"].Groupname, config["concurrence"].ServerParameter, config["concurrence"].Targetdbschema, config["concurrence"].Amount, c)
		c = c + offset - (c % offset)
	}
}
