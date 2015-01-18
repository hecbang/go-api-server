package testing

import (
	"const/path"
	"fmt"
	"libraries/common"
	"log"
	"net/http"
)

func AppConcurrence() {
	jq := common.NewJsonQuery(path.CONFIG_PATH + "testing" + path.DS + "app_concurrence.json")

	db := common.NewMySqlInstance("testdata")

	name := jq.String("group", "name")
	parameter := jq.String("group", "parameter")

	data := map[string]interface{}{
		"Name":              name,
		"SettingParameters": parameter,
		"LogTime":           common.Date("Y-m-d H:i:s"),
	}

	schemaName := jq.String("schema", "name")
	amount := jq.Int("schema", "amount")
	start := jq.Int("schema", "start")
	offset := jq.Int("schema", "offset")
	max := jq.Int("schema", "max")

	//处理掉重复的记录
	sql := "select app.Id, app.GroupId from app inner join app_group on app.GroupId=app_group.Id where app.Name=? and app_group.Name=?"
	result, err := db.GetRow(sql, schemaName, name)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if !common.Empty(result) {
		db.Delete("app_group", map[string]interface{}{"Id": result["GroupId"]})
		db.Delete("app", map[string]interface{}{"Id": result["Id"]})
	}

	lastid, err := db.Insert("app_group", data)
	if err != nil {
		panic(err.Error())
	}

	c := start
	fmt.Println("max concurrence is ", max)
	for c < max {
		fmt.Println("Now, concurrence = ", c)
		appconcurrence(lastid, schemaName, amount, c)
		c = c + offset - (c % offset)
	}
}

//app并发测试
//n 测试总次数
//c 并发量
func appconcurrence(groupid int64, schemaname string, n int, c int) {
	if c > n {
		panic("error: c>n")
	}
	timer := common.NewTimer()
	timer.Start()

	//每个线程执行多少次
	segs := n / c

	var chs []chan int = make([]chan int, c)
	var cycleN int
	for i := 0; i < c; i++ {
		cycleN = segs
		if i == c-1 {
			cycleN += n % c
		}
		chs[i] = make(chan int)
		go func(cycleN int, ch chan int) {
			for cyc := 0; cyc < cycleN; cyc++ {
				resp, err := http.Get("http://localhost:8080")
				if err != nil {
					log.Fatalln(err.Error())
				}
				resp.Body.Close()
			}
			ch <- 1
		}(cycleN, chs[i])
	}

	fmt.Println("wait...")
	for _, ch := range chs {
		<-ch
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
	tdb.Insert("app", result)
}
