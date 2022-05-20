package main

import (
	etcd "region/etcdManager"
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
	go etcd.ServiceRegister(global.Region)
	// 注册rpc服务
	rpc.RpcM2R, _ = rpc.DialReportService("tcp", global.MasterIP+":"+config.Configs.Rpc_M2R_port)

	// 向master报告本地的表
	regionWorker.SendNewTables("./Tables")
	// 发布rpc服务
	go rpc.StartCliService(config.Configs.Rpc_R2C_port)
	go rpc.StartGossipService(config.Configs.Rpc_R2R_port)
	for {

	}
}
