package parser

import (
	"miniSQL/src/Interpreter/types"
	"strconv"
	"strings"
	"testing"
)

var sql_strings = []string{
	"create table cxz(" +
		"afsdfsad int unique," +
		"what char(30) not null," +
		"primary key (what)" +
		");",
	"select a,b,c,d,e,f,g from cxz where a=123 and b=456 or c=234;",
}

var sql_CreateTable = "create table cxz(" +
	"afsdfsad int unique," +
	"what char(30) not null," +
	"primary key (what)" +
	");"
var sql_CreateDatabase = "create database cxz"

func TestCreateTable(t *testing.T) {
	ch := make(chan types.DStatements)
	go func() {
		io := strings.NewReader(sql_CreateTable) // 组装io
		Parse(io, ch)
	}()

	data := <-ch
	cts := data.(types.CreateTableStatement)
	t.Log(data.GetOperationName())
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
	t.Errorf("\n")
}

func TestCreateDatabaseStatement(t *testing.T) {
	ch := make(chan types.DStatements)
	go func() {
		io := strings.NewReader(sql_CreateDatabase) // 组装io
		Parse(io, ch)
	}()

	data := <-ch
	cts := data.(types.CreateDatabaseStatement)
	t.Log(data.GetOperationName())
	t.Log("要创建的数据库为:" + cts.DatabaseId)
	t.Errorf("\n")
}
