package main

import (
	"fmt"
	etcd "region/etcdManager"
	miniSQL "region/miniSQL"
	regionWorker "region/region"
	rpc "region/rpcManager"
	config "region/utils/ConfigSystem"
	mylog "region/utils/LogSystem"
	"region/utils/global"
)

func main() {
	// 完成初始化的准备工作
	mylog.LogInputChan = mylog.LogStart()
	config.BuildConfig()
	global.Region = etcd.Init()
	global.MasterIP = etcd.GetMasterIP(global.Region)
	global.AsyncLogSQLChan = regionWorker.StartAsyncCopy() // 开启全局的异步备份管道
	global.TableMap = make(map[string]*global.TableMeta)
	go etcd.ServiceRegister(global.Region)
	// 开启本地的SQL服务
	global.SQLInput = make(chan string, 10)
	global.SQLOutput = make(chan string, 10)
	go miniSQL.Start(global.SQLInput, global.SQLOutput)

	// buildSQL()

	global.SQLInput <- "use database aaa;"
	fmt.Println(1)
	res := <-global.SQLOutput
	fmt.Println(2)
	fmt.Println(res)

	// 注册rpc服务
	rpc.RpcM2R, _ = rpc.DialReportService("tcp", global.MasterIP+":"+config.Configs.Rpc_M2R_port)
	// 向master报告本地的表
	regionWorker.SendLocalTables(config.Configs.Minisql_table_store)
	// 发布rpc服务
	go rpc.StartCliService(config.Configs.Rpc_R2C_port)
	go rpc.StartGossipService(config.Configs.Rpc_R2R_port)

	for {

	}
}

// func buildSQL() {
// 	// global.SQLInput <- "create database aaa;"
// 	// res := <-global.SQLOutput
// 	// fmt.Println(res)
// 	// time.Sleep(3 * time.Second)
// 	fmt.Println(0)
// 	fmt.Println("管道的内容数量为:" + strconv.Itoa(len(global.SQLInput)))
// 	_, isOpen := <-global.SQLInput
// 	if !isOpen {
// 		fmt.Println("管道已被关闭")
// 	}
// 	global.SQLInput <- "use database aaa;"
// 	fmt.Println(1)
// 	res := <-global.SQLOutput
// 	fmt.Println(2)
// 	fmt.Println(res)
// }
