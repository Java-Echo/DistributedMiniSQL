package main

import (
	"fmt"
	miniSQL "region/miniSQL"
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
	// global.SQLInput <- "create database aaa;"
	// res := <-global.SQLOutput
	// fmt.Println(res)
	global.SQLInput <- "use database aaa;"
	res := <-global.SQLOutput
	fmt.Println(res)
	global.SQLInput <- "create table ttt(id int);"
	res = <-global.SQLOutput
	fmt.Println(res)
	t.Error("终止")
}
