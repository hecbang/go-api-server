package controllers

import (
	"const/path"
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

	schemaName, err := jq.String("schema", "name")
	if err != nil {
		log.Fatalln(err.Error())
	}

	targetSchemaDb, err := jq.String("schema", "db")
	if err != nil {
		log.Fatalln(err.Error())
	}

	amount, err := jq.Int("schema", "amount")
	if err != nil {
		log.Fatalln(err.Error())
	}

	start, err := jq.Int("schema", "start")
	if err != nil {
		log.Fatalln(err.Error())
	}

	offset, err := jq.Int("schema", "offset")
	if err != nil {
		log.Fatalln(err.Error())
	}

	max, err := jq.Int("schema", "max")
	if err != nil {
		log.Fatalln(err.Error())
	}

	c := start
	for c < max {
		testing.Concurrence(lastid, schemaName, targetSchemaDb, amount, c)
		c = c + offset - (c % offset)
	}
}
