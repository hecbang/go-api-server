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

	//处理掉重复的记录
	sql := "select db.Id, db.GroupId from db inner join db_group on db.GroupId=db_group.Id where db.Name=? and db_group.Name=?"
	result, err := db.GetRow(sql, schemaName, name)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if !common.Empty(result) {
		db.Delete("db_group", map[string]interface{}{"Id": result["GroupId"]})
		db.Delete("db", map[string]interface{}{"Id": result["Id"]})
	}

	lastid, err := db.Insert("db_group", data)

	c := start
	fmt.Println("max concurrence is ", max)
	for c < max {
		fmt.Println("Now, concurrence = ", c)
		testing.Concurrence(lastid, schemaName, targetSchemaDb, amount, c)
		c = c + offset - (c % offset)
	}
}
