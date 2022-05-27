package rpc

import (
	"fmt"
	"io"
	"os"
	"region/utils/global"
)

/*--------------master节点的操作--------------*/
// ToDo:查询语句
func MasterSQLSelect(SQL SQLRst) (string, bool) {
	global.SQLInput <- SQL.SQL
	res := <-global.SQLOutput
	return res, true
}

// ToDo:对部分表项的修改
func MasterSQLChange(SQL SQLRst) (string, bool) {
	global.SQLInput <- SQL.SQL
	res := <-global.SQLOutput
	fmt.Println(res)
	// SQL执行成功，尝试写日志
	if res == "Succeed" {
		writeSQLLog(SQL)
	}
	return res, true
}

// ToDo:对数据表的增加
func MasterSQLTableCreate(SQL SQLRst) (string, bool) {
	global.SQLInput <- SQL.SQL
	res := <-global.SQLOutput
	fmt.Println(res)
	// SQL执行成功，尝试写日志
	if res == "Succeed" {
		writeSQLLog(SQL)
	}
	return res, true
}

// ToDo:对数据表的删除
func MasterSQLTableDelete(SQL SQLRst) (string, bool) {
	global.SQLInput <- SQL.SQL
	res := <-global.SQLOutput
	return res, true
}

/*--------------常规的SQL操作--------------*/
// 所有的SQL语句都能执行，适用于日志的读取
func NormalSQL(sql string) string {
	fmt.Println("执行SQL语句:" + sql)
	global.SQLInput <- sql
	res := <-global.SQLOutput
	return res
}

// =======SQL的日志实现=======
func writeSQLLog(SQL SQLRst) {
	logName := SQL.Table + "_log"
	var logFile *os.File
	var err error
	if _, err := os.Stat(logName); os.IsNotExist(err) {
		// 这里理论上不应该创建一个文件，只有在create_table的时候才应该主动创建文件
		logFile, err = os.Create(logName) //创建文件
	} else {
		logFile, err = os.OpenFile(logName, os.O_APPEND|os.O_WRONLY, os.ModeAppend) //打开文件
		fmt.Println("文件存在")
	}
	defer logFile.Close()

	if SQL.SQLtype == "insert" || SQL.SQLtype == "update" || SQL.SQLtype == "delete" {
		// 写入日志
		_, err = io.WriteString(logFile, SQL.SQL+"\n")
		if err != nil {
			panic(err)
		}
	}

	if SQL.SQLtype == "create_table" {
		logFile, err = os.Create(logName) //创建文件
		// 会优先清除本地的表
		_, err = io.WriteString(logFile, "drop table "+SQL.Table+";"+"\n")
		_, err = io.WriteString(logFile, SQL.SQL+"\n")
	}
}
