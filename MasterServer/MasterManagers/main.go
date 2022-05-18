package main

import (
	"fmt"
	etcd "master/etcdManager"
	rpc "master/rpcManager"
	config "master/utils/ConfigSystem"
	mylog "master/utils/LogSystem"
	"master/utils/global"
)

func main() {
	fmt.Println("test")
	config.BuildConfig()
	global.Master = etcd.Init()
	go rpc.StartReportService()
	mylog.LogInputChan = mylog.LogStart()
	global.RegionMap = make(map[string]global.RegionMeta)
	global.TableMap = make(map[string]global.TableMeta)
	go etcd.RegisterWatcher(global.Master, config.Configs.Etcd_region_register_catalog)
	for {

	}
}
