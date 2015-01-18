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

//tp 支持两种取值 app, db
func (this *Testing) Concurrence(tp string) {
	jq := common.NewJsonQuery(path.CONFIG_PATH + "testing" + path.DS + "db_concurrence.json")

	db := common.NewMySqlInstance("testdata")

	name := jq.String("group", "name")
	parameter := jq.String("group", "parameter")

	data := map[string]interface{}{
		"Name":              name,
		"SettingParameters": parameter,
		"LogTime":           common.Date("Y-m-d H:i:s"),
	}

	schemaName := jq.String("schema", "name")
	targetSchemaDb := jq.String("schema", "db")
	amount := jq.Int("schema", "amount")
	start := jq.Int("schema", "start")
	offset := jq.Int("schema", "offset")
	max := jq.Int("schema", "max")

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
		testing.DbConcurrence(lastid, schemaName, targetSchemaDb, amount, c)
		c = c + offset - (c % offset)
	}
}
