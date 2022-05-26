package main

import (
	"fmt"
	"strconv"
	"testing"
)

func Test_parser(t *testing.T) {
	testSQL := []string{
		"select * from student2 where (((213!=id or 70<score) and 'name123'<name) or 99.0!=score) and name>'name9998';",
		"insert into student2 values(10866666666,'王振阳',100);",
		"update student2 set name='陈旭征' where name='王振阳';",
		"delete from student2 where name='name97996';",
		"drop table student2;",
		"create table ttt;",
	}
	for i, sql := range testSQL {
		s := parser(sql)
		fmt.Println("-------测试第" + strconv.Itoa(i) + "条语句-------")
		fmt.Println("SQL语句为:" + s.SQL)
		fmt.Println("SQL的类型为:" + s.SQLtype)
		fmt.Println("SQL执行的表为:" + s.Table)
	}
	// s := parser("select * from aaa;")
	// fmt.Println(s.SQL)
	// fmt.Println(s.SQLtype)
	// fmt.Println(s.Table)
	t.Error("终止")
}
