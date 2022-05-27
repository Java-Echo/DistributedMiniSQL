package main

import (
	"fmt"
	miniSQL "region/miniSQL"
	rpc "region/rpcManager"
	config "region/utils/ConfigSystem"
	mylog "region/utils/LogSystem"
	"region/utils/global"
	"testing"
)

func Test_main(t *testing.T) {
	mylog.LogInputChan = mylog.LogStart()
	config.BuildConfig()

	global.SQLInput = make(chan string)
	global.SQLOutput = make(chan string)
	go miniSQL.Start(global.SQLInput, global.SQLOutput)

	// 创建数据库
	global.SQLInput <- "create database aaa;"
	res := <-global.SQLOutput
	fmt.Println(res)
	// 使用数据库
	global.SQLInput <- "use database aaa;"
	res = <-global.SQLOutput
	fmt.Println(res)

	// 创建数据表
	// tableName := "ttt"
	// tableName := "bbb"
	tableName := "wrd"
	rpc.MasterSQLTableCreate(rpc.SQLRst{SQLtype: "create_table", SQL: "create table " + tableName + "(id int);", Table: tableName})
	// 插入数据
	rpc.MasterSQLChange(rpc.SQLRst{SQLtype: "insert", SQL: "insert into " + tableName + " values(1);", Table: tableName})
	rpc.MasterSQLChange(rpc.SQLRst{SQLtype: "insert", SQL: "insert into " + tableName + " values(2);", Table: tableName})
	rpc.MasterSQLChange(rpc.SQLRst{SQLtype: "insert", SQL: "insert into " + tableName + " values(3);", Table: tableName})

	// 尝试建立查询
	sqlRes, _ := rpc.MasterSQLSelect(rpc.SQLRst{SQLtype: "select", SQL: "select * from " + tableName + ";"})
	fmt.Println("sql的返回结果为:" + sqlRes)

	global.SQLInput <- "quit;"
	// res = <-global.SQLOutput

	t.Error("终止")
}
