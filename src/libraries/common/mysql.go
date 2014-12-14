package common

import (
	"const/path"
	"database/sql"
	//"fmt"
	_ "github.com/go-mysql-driver"
	"log"
	"regexp"
	"strings"
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

	var key string = BuildKeyMd5(schema, conntype)
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
	var key string = BuildKeyMd5(schema, contype)
	if contype == "Master" {
		dataSourceName = conf.Master
	} else {
		dataSourceName = conf.Slave[Rand(0, len(conf.Slave)-1)]
	}

	//开始连接DB
	dbinit, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatalln(err.Error())
	}

	//将DB连接放入一个全局变量中
	db[key] = dbinit
}

//根据所提供的sql语句获取数据列表
func (this *MySql) GetAll(sql string, args ...interface{}) ([]map[string]string, error) {
	this.checkSQL(sql)
	rows, err := this.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return this.fetch(rows)
}

//根据所提供的sql语句获取数据列表
func (this *MySql) GetAllArray(sql string, args ...interface{}) ([]map[int]string, error) {
	this.checkSQL(sql)
	rows, err := this.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return this.fetchArray(rows)
}

//根据SQL获取一行数据
func (this *MySql) GetRow(sql string, args ...interface{}) (map[string]string, error) {
	this.checkSQL(sql)
	if !strings.Contains(strings.ToLower(sql), "limit") {
		sql += " LIMIT 1"
	}
	rows, err := this.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result, err := this.fetch(rows)
	if err != nil {
		return nil, err
	}

	retval := make(map[string]string, 0)
	if len(result) > 0 {
		retval = result[0]
	}
	return retval, nil
}

//根据SQL获取一行数据
func (this *MySql) GetRowArray(sql string, args ...interface{}) (map[int]string, error) {
	this.checkSQL(sql)
	if !strings.Contains(strings.ToLower(sql), "limit") {
		sql += " LIMIT 1"
	}
	rows, err := this.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result, err := this.fetchArray(rows)
	if err != nil {
		return nil, err
	}

	retval := make(map[int]string, 0)
	if len(result) > 0 {
		retval = result[0]
	}
	return retval, nil
}

//获了结果集中的第一行第一列元素值
func (this *MySql) GetOne(sql string, args ...interface{}) (string, error) {
	result, err := this.GetAllArray(sql, args...)
	if err != nil {
		return "", err
	}
	if len(result) > 0 {
		return result[0][0], nil
	}
	return "", nil
}

//根据表名、指定字段、条件获取数据列表
//参数1: 表名
//参数2: 字段列表[]string
//参数3: 查询条件map[string]string
//参数4: 排序条件string，格式如："Id DESC"
//参数5: 取记录限制string，格式如"10, 20" 或 "5"
func (this *MySql) GetList(table string, args ...interface{}) ([]map[string]string, error) {
	var fields []string = make([]string, 0)
	var conditions map[string]string = make(map[string]string)
	var orderby string
	var limit string
	if len(args) > 0 {
		for i, v := range args {
			switch i {
			case 0:
				fields = v.([]string)
			case 1:
				conditions = v.(map[string]string)
			case 2:
				orderby = v.(string)
			case 3:
				limit = v.(string)
			}
		}
	}
	//开始组装SQL语句
	querysql, arguments := this.buildSql(table, fields, conditions, orderby, limit)
	return this.GetAll(querysql, arguments...)
}

//根据表名、指定字段、条件获取一条数据记录
//参数1: 表名
//参数2: 字段列表[]string
//参数3: 查询条件map[string]string
//参数4: 排序条件string，格式如："Id DESC"
func (this *MySql) GetDictionary(table string, args ...interface{}) (map[string]string, error) {
	var fields []string = make([]string, 0)
	var conditions map[string]string = make(map[string]string)
	var orderby string
	if len(args) > 0 {
		for i, v := range args {
			switch i {
			case 0:
				fields = v.([]string)
			case 1:
				conditions = v.(map[string]string)
			case 2:
				orderby = v.(string)
			}
		}
	}
	//开始组装SQL语句
	querysql, arguments := this.buildSql(table, fields, conditions, orderby, "1")
	return this.GetRow(querysql, arguments...)
}

func (this *MySql) buildSql(table string, fields []string, conditions map[string]string, orderby, limit string) (string, []interface{}) {
	var arguments []interface{} = make([]interface{}, 0)
	//拼接field部分
	var fieldstr string = "*"
	if !Empty(fields) {
		fieldstr = strings.Join(fields, ", ")
	}

	//拼接condition条件部分
	var condlist []string = make([]string, 0)
	for k, v := range conditions {
		vlist := StringToList(v)
		if len(vlist) == 1 {
			arguments = append(arguments, vlist[0])
			condlist = append(condlist, k+"=?")
		} else if len(vlist) > 1 {
			placeholders := make([]string, 0)
			for _, val := range vlist {
				arguments = append(arguments, val)
				placeholders = append(placeholders, "?")
			}
			condlist = append(condlist, k+" IN("+strings.Join(placeholders, ",")+")")
		}
	}
	var condstr string = ""
	if !Empty(condlist) {
		condstr = " WHERE " + strings.Join(condlist, " AND ")
	}

	//order by
	if !Empty(orderby) {
		orderby = " ORDER BY " + orderby
	}

	//limit
	if !Empty(limit) {
		limit = " LIMIT " + limit
	}
	querysql := "SELECT " + fieldstr + " FROM " + table + condstr + orderby + limit
	this.checkSQL(querysql)
	return querysql, arguments
}

//根据db.Query的查询结果，组装成一个关联key的数据集，数据类型[]map[string]string
func (this *MySql) fetch(rows *sql.Rows) ([]map[string]string, error) {
	result := make([]map[string]string, 0)
	columns, err := rows.Columns()
	if err != nil {
		//an error occurred
		return nil, err
	}

	rawBytes := make([]sql.RawBytes, len(columns))

	//rows.Scan wants '[]interface{}' as an argument, so we must copy
	//the references into such a slice
	scanArgs := make([]interface{}, len(columns))

	for i := range rawBytes {
		scanArgs[i] = &rawBytes[i]
	}

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		var val string
		item := make(map[string]string)
		for i, col := range rawBytes {
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

//根据db.Query的查询结果，组装成一个array的数据集，数据类型[]map[int]string
func (this *MySql) fetchArray(rows *sql.Rows) ([]map[int]string, error) {
	result := make([]map[int]string, 0)
	columns, err := rows.Columns()
	if err != nil {
		//an error occurred
		return nil, err
	}

	rawBytes := make([]sql.RawBytes, len(columns))

	//rows.Scan wants '[]interface{}' as an argument, so we must copy
	//the references into such a slice
	scanArgs := make([]interface{}, len(columns))

	for i := range rawBytes {
		scanArgs[i] = &rawBytes[i]
	}

	item := make(map[int]string)
	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		var val string
		for i, col := range rawBytes {
			if col == nil {
				val = ""
			} else {
				val = string(col)
			}
			item[i] = val
		}
		result = append(result, item)
	}
	return result, nil
}

//向前一个表中写入一条记录，如果写入成功，则返回其主键ID值
func (this *MySql) Insert(table string, data map[string]interface{}) (int64, error) {
	columns, err := this.GetTableColumns(table)
	if err != nil {
		return 0, err
	}
	log.Println(columns)
	//TODO

	return int64(64), nil
}

//获取一个表的列信息，字段名称为key，值为字段信息k-v
func (this *MySql) GetTableColumns(table string) (map[string]map[string]string, error) {
	list, err := this.GetAll("DESC " + table)
	if err != nil {
		return nil, err
	}
	retval := make(map[string]map[string]string, 0)
	for _, item := range list {
		retval[item["Field"]] = item
	}
	return retval, nil
}

//保证修改、写入类的操作不在slave上执行
func (this *MySql) checkSQL(sql string) {
	if this.optype == "Slave" {
		sql = strings.TrimSpace(sql)
		exp := regexp.MustCompile(`^(?i:insert|update|delete|alter|truncate|drop)`)
		if exp.MatchString(sql) {
			panic("insert|update|delete|alter|truncate|drop operation is not allowed in slave.")
		}
	}
}
