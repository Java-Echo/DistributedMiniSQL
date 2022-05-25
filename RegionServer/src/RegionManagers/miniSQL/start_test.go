package main

import (
	"testing"
	"time"
)

func Test_start(t *testing.T) {
	in := make(chan string, 1)
	out := make(chan string, 1)
	go start(in, out)
	in <- "use database aaa;"
	<-out
	in <- "create table ttt( id int);"
	CTres := <-out
	t.Log("创建表：" + CTres)
	in <- "select * from ttt;"
	str := <-out
	t.Log("查询的结果为" + str)
	time.Sleep(5 * time.Second)
	t.Error("终止")
}
