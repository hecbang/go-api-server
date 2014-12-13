package common

import (
	"const/path"
	"database/sql"
	//"fmt"
	_ "github.com/go-mysql-driver"
	"log"
)

var db map[string]*sql.DB

func init() {
	db = make(map[string]*sql.DB)
}

type MySql struct {
	db     *sql.DB
	optype string
}

//创建一个默认的mysql操作实例
func NewMySql() *MySql {
	return NewMySqlInstance("default", "Master")
}

//创建一个默认的mysql查询实例
func NewQuery() *MySql {
	return NewMySqlInstance("default", "Slave")
}

//实例化一个mysql实例
//@param string schema 数据库连接方案
//@param string conntype 数据库连接类型，范围: Master, Slave
func NewMySqlInstance(schema string, conntype string) *MySql {
	if !InList(conntype, []string{"Master", "Slave"}) {
		panic("function common.NewSqlInstance's second argument must be Master or Slave.")
	}

	var key string = schema + conntype
	_, ok := db[key]
	if !ok {
		//建立一个新连接到mysql
		connect(schema, conntype)
	}
	return &MySql{db: db[key], optype: conntype}
}

//建立数据库连接
//@param string schema 连接DB方案
//@param string conntype 连接类型，是分Master和Slave类型
func connect(schema string, contype string) {
	type item struct {
		Master string
		Slave  []string
	}
	var v map[string]item

	//获取DB连接配置文件
	if err := LoadJson(path.CONFIG_PATH+"db.json", &v); err != nil {
		log.Fatalln(err.Error())
	}
	conf, ok := v[schema]
	if !ok {
		log.Fatalln("Database configuration file error. Lost schema[" + schema + "] node.")
	}
	var dataSourceName string
	var key string
	if contype == "Master" {
		dataSourceName = conf.Master
		key = schema + "Master"
	} else {
		dataSourceName = conf.Slave[Rand(0, len(conf.Slave)-1)]
		key = schema + "Slave"
	}

	//开始连接DB
	dbinit, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatalln(err.Error())
	}

	//将DB连接放入一个全局变量中
	db[key] = dbinit
}
