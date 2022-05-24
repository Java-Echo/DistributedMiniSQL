package main

import (
	etcd "master/etcdManager"
	rpc "master/rpcManager"
	config "master/utils/ConfigSystem"
	mylog "master/utils/LogSystem"
	"master/utils/global"
)

func main() {
	mylog.LogInputChan = mylog.LogStart()
	config.BuildConfig()
	global.Master = etcd.Init()
	global.RegionMap = make(map[string]*global.RegionMeta)
	global.TableMap = make(map[string]*global.TableMeta)
	// 发布rpc服务
	go rpc.StartReportService(config.Configs.Rpc_m2r_port)
	go rpc.StartCliService(config.Configs.Rpc_m2c_port)
	// go etcd.RegisterWatcher(global.Master, config.Configs.Etcd_region_register_catalog)
	go etcd.RegisterWatcherWithWorker(global.Master, config.Configs.Etcd_region_register_catalog, &etcd.RegionRegisterWorker{})
	go etcd.RegisterWatcherWithWorker(global.Master, config.Configs.Etcd_region_stepout_catalog, &etcd.RegionStepOutWorker{})
	for {

	}
}
