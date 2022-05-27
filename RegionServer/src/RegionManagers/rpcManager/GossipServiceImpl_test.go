package rpc

import (
	"fmt"
	"io/ioutil"
	"log"
	config "region/utils/ConfigSystem"
	"testing"
)

func TestGossipService_PassLog(t *testing.T) {
	client, err := DialGossipService("tcp", "localhost:"+config.Configs.Rpc_R2R_port)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	var reply PassLogRes
	request := PassLogRst{
		SqlType: "insert",
		Sql:     "inset xxx",
		Table:   "aaa",
	}
	err = client.PassLog(request, &reply)
	if err != nil {
		log.Fatal(err)
	}
	t.Log("成功执行了应该，我猜的")
	t.Error("终止")
}

// func TestGossipService_PassTable(t *testing.T) {
// 	testFile := "这是一段测试的代码"
// 	testByte := []byte(testFile)
// 	// ToDo:将其改为目标机器的ip地址
// 	client, err := DialGossipService("tcp", "localhost:"+config.Configs.Rpc_R2R_port)
// 	if err != nil {
// 		log.Fatal("dialing:", err)
// 	}
// 	var reply PassTableRes
// 	request := PassTableRst{
// 		Content:   testByte,
// 		TableName: "随便什么表格",
// 	}
// 	err = client.PassTable(request, &reply)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	t.Log("成功执行了应该，我猜的")
// 	t.Error("终止")
// }

func TestGossipService_PassTable(t *testing.T) {
	logFile, err := ioutil.ReadFile("./ttt_log")
	if err != nil {
		fmt.Println("read fail", err)
	}

	// ToDo:将其改为目标机器的ip地址
	client, err := DialGossipService("tcp", "10.162.19.119:"+config.Configs.Rpc_R2R_port)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	var reply PassTableRes
	request := PassTableRst{
		Content:   logFile,
		TableName: "ttt",
	}
	err = client.PassTable(request, &reply)
	if err != nil {
		log.Fatal(err)
	}
	t.Log("成功执行了应该，我猜的")
	t.Error("终止")
}

func TestGossipService_SyncSQL(t *testing.T) {
	client, err := DialGossipService("tcp", "10.162.19.119:"+config.Configs.Rpc_R2R_port)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	var reply SQLRes
	request := SQLRst{
		SQLtype: "insert",
		SQL:     "insert into ttt values(1);",
		Table:   "ttt",
	}
	err = client.SyncSQL(request, &reply)
}
