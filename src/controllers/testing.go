package controllers

import (
	"const/path"
	"fmt"
	"io/ioutil"
	"libraries/common"
	"log"
	"models/testing"
	"strconv"
	"strings"
)

type Testing struct {
}

//tp 支持两种取值 app, db
func (this *Testing) Concurrence(tp string) {
	if tp == "db" {
		testing.DatabaseConcurrence()
	} else {
		testing.AppConcurrence()
	}
}

func (this *Testing) Init() {
	jq := common.NewJsonQuery(path.CONFIG_PATH + "testing" + path.DS + "db_concurrence.json")
	bytes, err := common.Read("slap_raw.sh")
	if err != nil {
		log.Fatalln(err.Error())
	}
	content := string(bytes)
	is_query := jq.String("is_query")
	host := jq.String("db", "host")
	port := jq.String("db", "port")
	user := jq.String("db", "user")
	password := jq.String("db", "password")
	dbname := jq.String("db", "dbname")
	amount := jq.String("schema", "amount")
	concurrency := jq.String("schema", "concurrency")
	content = strings.Replace(content, "#is_query#", is_query, -1)
	content = strings.Replace(content, "#host#", host, -1)
	content = strings.Replace(content, "#port#", port, -1)
	content = strings.Replace(content, "#user#", user, -1)
	content = strings.Replace(content, "#password#", password, -1)
	content = strings.Replace(content, "#dbname#", dbname, -1)
	content = strings.Replace(content, "#concurrency#", concurrency, -1)
	content = strings.Replace(content, "#amount#", amount, -1)
	ioutil.WriteFile("slap.sh", []byte(content), 0755)
}

func (this *Testing) Parser() {
	jq := common.NewJsonQuery(path.CONFIG_PATH + "testing" + path.DS + "db_concurrence.json")

	name := jq.String("group", "name")
	parameter := jq.String("group", "parameter")

	//名称格式
	schemaName := jq.String("schema", "name")

	//测试执行次数	和并发数
	amount := jq.Int("schema", "amount")
	concurrency := jq.String("schema", "concurrency")
	concurrencyList := common.StringToList(concurrency)

	is_query := jq.String("is_query")

	results, err := common.ReadLine("result.info")
	if err != nil {
		log.Fatalln(err.Error())
	}

	size := len(results)
	if size > 0 && len(results) != len(concurrencyList) {
		fmt.Println("Error, the number if concurrency is not matched the number of results.")
		return
	}

	//写入数据
	db := common.NewMySqlInstance("testdata")
	//处理掉重复的记录
	sql := "select db.GroupId from db inner join db_group on db.GroupId=db_group.Id where db.Name=? and db_group.Name=?"
	result, err := db.GetRow(sql, schemaName, name)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if !common.Empty(result) {
		db.Delete("db_group", map[string]interface{}{"Id": result["GroupId"]})
		db.Delete("db", map[string]interface{}{"GroupId": result["GroupId"]})
	}

	//写入分组数据
	data := map[string]interface{}{
		"Name":              name,
		"SettingParameters": parameter,
		"LogTime":           common.Date("Y-m-d H:i:s"),
	}
	groupId, err := db.Insert("db_group", data)
	if err != nil {
		panic(err.Error())
	}

	for i := 0; i < size; i++ {
		c, err := strconv.Atoi(concurrencyList[i])
		if err != nil {
			panic(err.Error())
		}
		elapse, err1 := strconv.ParseFloat(strings.TrimSpace(results[i]), 64)
		if err1 != nil {
			panic(err1.Error())
		}
		n := amount
		//convert to ms
		elapse = elapse * 1000

		if is_query == "1" {
			//查询的情况，小于50的并发，约定取并发数的50倍
			if c < 50 {
				n = c * 50
			}
		}
		result := map[string]interface{}{
			"GroupId":     groupId,
			"Name":        schemaName,
			"Total":       n,
			"Concurrence": c,
			"ElapseTime":  elapse,
			"QPS":         common.Round(float64(n)*1000/float64(elapse), 2), //每秒处理请求数
			"TPQ":         common.Round(float64(elapse)/float64(n), 2),      //平均每个请求用时多少ms
			"LogTime":     common.Date("Y-m-d H:i:s"),
		}
		db.Insert("db", result)
	}
}
