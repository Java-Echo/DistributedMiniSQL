package rpc

import (
	"fmt"
	"region/miniSQL"
	config "region/utils/ConfigSystem"
	mylog "region/utils/LogSystem"
	"region/utils/global"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	mylog.LogInputChan = mylog.LogStart()
	config.BuildConfig()
	global.SQLInput = make(chan string)
	global.SQLOutput = make(chan string)
	go miniSQL.Start(global.SQLInput, global.SQLOutput)

	// global.SQLInput <- "create database aaa;"
	// res := <-global.SQLOutput
	// fmt.Println(res)
	global.SQLInput <- "use database aaa;"
	res := <-global.SQLOutput
	fmt.Println(res)
	fmt.Println("初始化完成")

	m.Run()
	global.SQLInput <- "quit;"
}

func TestMasterSQLSelect(t *testing.T) {
	// sql := "select * from ttt;"
	sql := SQLRst{
		SQLtype: "select",
		SQL:     "select * from ttt;",
		Table:   "ttt",
	}
	res, ok := MasterSQLSelect(sql)
	if ok {
		fmt.Println("查询结果为:" + res)
	}
	time.Sleep(1 * time.Second)
	t.Error("终止")
}

func TestMasterSQLTableCreate(t *testing.T) {
	sql := SQLRst{
		SQLtype: "create_table",
		SQL:     "create table ttt(id int);",
		Table:   "ttt",
	}
	res, ok := MasterSQLTableCreate(sql)
	if ok {
		fmt.Println("表的创建结果为:" + res)
	}
	time.Sleep(1 * time.Second)
	t.Error("终止")
}

func TestMasterSQLChange(t *testing.T) {
	// sql := "insert into ttt values(3);"
	sql := SQLRst{
		SQLtype: "insert",
		SQL:     "insert into ttt values(3);",
		Table:   "ttt",
	}
	res, ok := MasterSQLTableCreate(sql)
	if ok {
		fmt.Println("表项插入的结果为:" + res)
	}
	time.Sleep(1 * time.Second)
	t.Error("终止")
}
