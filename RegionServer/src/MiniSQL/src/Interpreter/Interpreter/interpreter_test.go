package interpreter

import (
	"miniSQL/src/Interpreter/types"
	"strconv"
	"testing"
	"time"
)

var sql_strings = []string{
	"create table cxz(" +
		"afsdfsad int unique," +
		"what char(30) not null," +
		"primary key (what)" +
		");",
	"select a,b,c,d,e,f,g from cxz where a=123 and b=456 or c=234;",
	// "insert into student2 values(1080100001,'name1',99);",
	"create index name_index on student2 ( score);",
	"drop index name_index on student2;",
	"drop table student2;",
	"delete from student2 where name='name97996';",
	"update student2 set name='aaa' where name='bbb';",
	"use database wzy;",
	// "update student2 set name='陈旭征' where name='王振阳'; ",
}

func TestParser(t *testing.T) {
	input := make(chan string)
	output := CreateInterpreter(input)
	go func() {
		for _, sql := range sql_strings {
			input <- sql
		}
	}()

	go func() {
		for {
			sqlStmt := <-output
			t.Log(sqlStmt.GetOperationName())
		}
	}()

	time.Sleep(2 * time.Second)
	t.Errorf("不要慌，我只是想要一个log而已\n")
}

var sql_CreateTable = "create table cxz(" +
	"afsdfsad int unique," +
	"what char(30) not null," +
	"primary key (what)" +
	");"
var sql_CreateDatabase = "create database cxz"
var sql_CreateIndex = "create index name_index on student2 ( score);"
var sql_Select = "select * from student2 where name='name97996';"

func TestCreateTable(t *testing.T) {
	input := make(chan string)
	output := CreateInterpreter(input)
	go func() {
		input <- sql_CreateTable
	}()

	go func() {
		sqlStmt := <-output
		cts := sqlStmt.(types.CreateTableStatement)
		t.Log(sqlStmt.GetOperationName())
		t.Log("要创建的表为:" + cts.TableName)

		keys := ""
		for key, attr := range cts.ColumnsMap {
			attrDetails := ""
			attrDetails += "(类型为:" + strconv.Itoa(attr.Type.TypeTag)
			if attr.Unique {
				attrDetails += ", 唯一"
			}
			if attr.NotNull {
				attrDetails += ", 非空"
			}
			attrDetails += ")"
			keys += " " + key + attrDetails
		}
		t.Log("加入的键为:" + keys)

		primaryKey := ""
		for _, key := range cts.PrimaryKeys {
			primaryKey += key.Name + " "
		}
		t.Log("主键为:" + primaryKey)
	}()

	time.Sleep(1 * time.Second)
	t.Errorf("\n")
}

func TestCreateDatabaseStatement(t *testing.T) {
	input := make(chan string)
	output := CreateInterpreter(input)
	go func() {
		input <- sql_CreateDatabase
	}()

	go func() {
		sqlStmt := <-output
		t.Log(sqlStmt.GetOperationName())
		cds := sqlStmt.(types.CreateDatabaseStatement)
		t.Log("要创建的数据库为:" + cds.DatabaseId)
	}()

	time.Sleep(1 * time.Second)
	t.Errorf("\n")
}

func TestCreateIndexStatement(t *testing.T) {
	input := make(chan string)
	output := CreateInterpreter(input)
	go func() {
		input <- sql_CreateIndex
	}()

	go func() {
		sqlStmt := <-output
		t.Log(sqlStmt.GetOperationName())
		cis := sqlStmt.(types.CreateIndexStatement)
		res := cis.IndexName
		if cis.Unique {
			res += "(唯一)"
		}
		res += "作用在表" + cis.TableName + "上  "
		res += "作用的键为:"
		for _, key := range cis.Keys {
			res += " " + key.Name
		}
		t.Log(res)
	}()

	time.Sleep(1 * time.Second)
	t.Errorf("\n")
}

func TestSelectStatement(t *testing.T) {
	input := make(chan string)
	output := CreateInterpreter(input)
	go func() {
		input <- sql_Select
	}()

	go func() {
		sqlStmt := <-output
		t.Log(sqlStmt.GetOperationName())
		// ss := sqlStmt.(types.SelectStatement)
		// ToDo: 完善select的测试样例

	}()

	time.Sleep(1 * time.Second)
	t.Errorf("\n")
}

// func Test...(t *testing.T) {
// 	input := make(chan string)
// 	output := CreateInterpreter(input)
// 	go func() {
// 		input <- sql_...
// 	}()

// 	go func() {
// 		sqlStmt := <-output
//		t.Log(sqlStmt.GetOperationName())
// 		cts := sqlStmt.(types....)

// 	}()

// 	time.Sleep(1 * time.Second)
// 	t.Errorf("\n")
// }
