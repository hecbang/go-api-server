package testing

import (
	"const/path"
	"fmt"
	"libraries/common"
	"log"
	"strings"
)

func DatabaseConcurrence() {
	jq := common.NewJsonQuery(path.CONFIG_PATH + "testing" + path.DS + "db_concurrence.json")

	name := jq.String("group", "name")
	parameter := jq.String("group", "parameter")

	schemaNameFormat := jq.String("schema", "name")

	targetSchemaDb := jq.String("schema", "db")

	//测试重启执行次数
	times := jq.Int("schema", "times")
	amount := jq.Int("schema", "amount")
	start := jq.Int("schema", "start")
	offset := jq.Int("schema", "offset")
	max := jq.Int("schema", "max")

	var schemaName string
	for seq := 1; seq <= times; seq++ {
		//给测试方案名称加上顺号，标识不同的测试次数
		if times == 1 {
			schemaName = strings.Replace(schemaNameFormat, "(%d)", "", 1)
		} else {
			fmt.Println(">>>>>>database pressure test>>>>>>", seq)
			schemaName = fmt.Sprintf(schemaNameFormat, seq)
		}

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
		lastid, err := db.Insert("db_group", data)
		if err != nil {
			panic(err.Error())
		}

		c := start
		fmt.Println("max concurrence is ", max)
		for c < max {
			fmt.Println("Now, concurrence = ", c)
			dbconcurrence(lastid, schemaName, targetSchemaDb, amount, c)
			if c < 100 {
				c = c + 10 - (c % 10)
			} else {
				c = c + offset - (c % offset)
			}
		}
	}

}

//DB并发测试
//n 测试总次数
//c 并发量
func dbconcurrence(groupid int64, schemaname string, targetdbschema string, n int, c int) {
	if c > n {
		panic("error: c>n")
	}

	//先清空写入的目标数据库
	db := common.NewMySqlInstance(targetdbschema)
	_, err := db.UDExec("truncate table target")
	if err != nil {
		log.Fatalln(err.Error())
	}

	timer := common.NewTimer()
	timer.Start()

	//每个线程执行多少次
	segs := n / c

	var chs chan int = make(chan int, c)
	var cycleN int
	for i := 0; i < c; i++ {
		cycleN = segs
		if i == c-1 {
			cycleN += n % c
		}
		go func(cycleN int, chs chan int) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("%v", r)
					chs <- 1
				}
			}()
			db := common.NewMySqlInstance(targetdbschema)
			date := common.Date("Y-m-d H:i:s")
			for cyc := 0; cyc < cycleN; cyc++ {
				_, err := db.InsertExec("insert into target(Num, String, LogTime) values(?,?,?)", 98, "helloworld", date)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
			chs <- 1
		}(cycleN, chs)
	}

	fmt.Println("wait...")
	for i := 0; i < c; i++ {
		<-chs
	}

	elapse := timer.Elapse("ms")

	fmt.Println("finished once concurrence, elapse ", elapse, "ms")

	result := map[string]interface{}{
		"GroupId":     groupid,
		"Name":        schemaname,
		"Total":       n,
		"Concurrence": c,
		"ElapseTime":  elapse,
		"QPS":         common.Round(float64(n)/(float64(elapse)/1000), 2), //每秒处理请求数
		"TPQ":         common.Round(float64(elapse)/float64(n), 2),        //平均每个请求用时多少ms
		"LogTime":     common.Date("Y-m-d H:i:s"),
	}
	tdb := common.NewMySqlInstance("testdata")
	tdb.Insert("db", result)
}
