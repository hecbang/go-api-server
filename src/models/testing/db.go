package testing

import (
	"const/path"
	"fmt"
	"libraries/common"
	"log"
	"strings"
)

//数据库测试并发执行逻辑
func DatabaseConcurrence() {
	jq := common.NewJsonQuery(path.CONFIG_PATH + "testing" + path.DS + "db_concurrence.json")

	name := jq.String("group", "name")
	parameter := jq.String("group", "parameter")
	schemaNameFormat := jq.String("schema", "name")
	targetSchemaDb := jq.String("schema", "db")
	is_query := jq.Int("is_query")

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
			if is_query == 0 {
				dbconcurrence(lastid, schemaName, targetSchemaDb, amount, c)
			} else if is_query == 1 {
                queryconcurrence(lastid, schemaName, targetSchemaDb, amount, c)                
			} else {
                //example is_query=2
                //slow query
                adjust := c*32
                if adjust < amount {
                    fmt.Println("Adjust amount to ", adjust)
                    queryconcurrence(lastid, schemaName, targetSchemaDb, adjust, c)   
                } else {
                    queryconcurrence(lastid, schemaName, targetSchemaDb, amount, c)
                }
            }

			if c < 100 {
				c = c + 10 - (c % 10)
			} else {
				c = c + offset - (c % offset)
			}
		}
	}

}

//DB并发测试
//n 测试写入总次数
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

//DB并发测试
//n 测试查询总次数
//c 并发量
func queryconcurrence(groupid int64, schemaname string, targetdbschema string, n int, c int) {
	if c > n {
		panic("error: c>n")
	}

	//获取要查询的目标数据表和字段
	fmt.Println("Initialization data...")
	jq := common.NewJsonQuery(path.CONFIG_PATH + "testing" + path.DS + "db_concurrence.json")
	table := jq.String("query", "table")
	field := jq.String("query", "condition_field")	

	db := common.NewMySqlInstance(targetdbschema)
	keylist, err := db.GetAll("select distinct "+field+" from "+table+" limit ?", n)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if len(keylist) < n {
		log.Fatalln("Failed! Test data is less than ", n)
	}

	seeds := make([]string, n)
	for i, itm := range keylist {
		seeds[i] = itm[field]
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
		go func(table string, field string, seeds []string, page int, segs int, cycleN int, chs chan int) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("%v", r)
					chs <- 1
				}
			}()
			db := common.NewMySqlInstance(targetdbschema)
			//开始查询
			for k := 0; k < cycleN; k++ {
				key := segs * page
				//res, err := db.GetAll("select * from "+table+" where "+field+"=?", seeds[key])
                sql := "select count(EventTypeId) as cnt, EventTypeId from server_alarm_notice_item where CreateTime between '2014-06-01' and '2014-08-21' group by EventTypeId";
                _, err := db.GetAll(sql)
				if err != nil {
					log.Fatalln(err.Error())
				}
				/*if common.Empty(res) {
					log.Fatalln("No data fetched. " + field + "=" + seeds[key])
				}*/
				key++
			}

			chs <- 1
		}(table, field, seeds, i, segs, cycleN, chs)
	}

	fmt.Println("wait...")
	for i := 0; i < c; i++ {
		<-chs
	}

	elapse := timer.Elapse("ms")

	fmt.Println("finished once query concurrence, elapse ", elapse, "ms")

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
