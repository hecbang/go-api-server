package mysql

import (
	"database/sql"
	_ "github.com/go-mysql-driver"
	"log"
	"strconv"
	"strings"
	"time"
)

var db *sql.DB
var connTime int64 = 0
var waitTimeout int64 = 0

func init() {
	connect()
	checkConn()
}

//数据库连接操作
func connect() {
	initdb, err := sql.Open("mysql", "root:84373723@tcp(127.0.0.1:3306)/pharmacy?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}
	db = initdb
	connTime = time.Now().Unix()
	waitTimeout = getWaitTimeout()
}

//检查并推断数据库连接是否可用，如果不可用，则重新建立连接
func checkConn() {
	if waitTimeout == 0 {
		waitTimeout = getWaitTimeout()
	}
	if time.Now().Unix()-connTime+2 > waitTimeout {
		connect()
	}
	connTime = time.Now().Unix()
}

func getWaitTimeout() int64 {
	//定位数据库对空闲连接的等待时长
	var mysql MySql
	rows, err := db.Query("SHOW VARIABLES LIKE 'wait_timeout'")
	if err != nil {
		log.Fatal(err.Error())
	}
	result, err := mysql.Fetch(rows)
	if len(result) > 0 {
		time, err1 := strconv.Atoi(result[0]["Value"])
		if err1 != nil {
			log.Fatal(err1.Error())
		}
		return int64(time)
	} else {
		log.Fatal(err.Error())
	}
	return 0
}

type MySql struct {
}

//根据所提供的SQL语句获取数据列表
func (mysql *MySql) GetAll(sql string, args ...interface{}) ([]map[string]string, error) {
	checkConn()
	rows, err := db.Query(sql, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	return mysql.Fetch(rows)
}

//根据SQL语句获取一行数据
func (mysql *MySql) GetRow(sql string, args ...interface{}) (map[string]string, error) {
	retval := make(map[string]string, 0)
	if !strings.Contains(strings.ToLower(sql), "limit") {
		sql += " LIMIT 1"
	}
	rows, err := db.Query(sql, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	result, err := mysql.Fetch(rows)
	if err != nil {
		return retval, err
	}
	if len(result) > 0 {
		retval = result[0]
	}
	return retval, nil

}

//根据表名、指定字段、条件获取数据列表
func (mysql *MySql) GetList(table string, fields []string, conditions map[string]string) ([]map[string]string, error) {
	checkConn()
	//拼接field部分
	querysql := mysql.buildSql(table, fields, conditions)

	rows, err := db.Query(querysql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	return mysql.Fetch(rows)
}

func (mysql *MySql) Fetch(rows *sql.Rows) ([]map[string]string, error) {
	result := make([]map[string]string, 0)
	columns, err := rows.Columns()
	if err != nil {
		//data not found
		return result, nil
	}

	//make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	//rows.Scan wants '[]interface{}' as an argument, so we must copy the
	//references into such a slice
	scanArgs := make([]interface{}, len(columns))

	//the type '[]interface{}' references to '[]sql.RawBytes'
	for i := range values {
		scanArgs[i] = &values[i]
	}

	item := make(map[string]string)
	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Fatal(err)
		}

		var val string
		for i, col := range values {
			if col == nil {
				val = ""
			} else {
				val = string(col)
			}
			item[columns[i]] = val
		}
		result = append(result, item)
	}
	return result, nil
}

//根据表、字段、条件拼接SQL语句
func (mysql *MySql) buildSql(table string, fields []string, conditions map[string]string) string {
	//拼接field部分
	fieldstr := strings.Join(fields, ", ")

	//拼接condition条件部分
	var condlist []string
	for k, v := range conditions {
		condlist = append(condlist, k+"='"+v+"'")
	}
	condstr := strings.Join(condlist, " AND ")

	querysql := "SELECT " + fieldstr + " FROM " + table + " WHERE " + condstr
	return querysql
}
