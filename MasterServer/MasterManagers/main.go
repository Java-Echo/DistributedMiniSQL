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
	go rpc.StartReportService()
	global.RegionMap = make(map[string]global.RegionMeta)
	global.TableMap = make(map[string]global.TableMeta)
	go etcd.RegisterWatcher(global.Master, config.Configs.Etcd_region_register_catalog)
	for {

	}
}
