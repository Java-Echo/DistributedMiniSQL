package main

import (
	config "client/utils/ConfigSystem"
	mylog "client/utils/LogSystem"
	"client/utils/global"
)

func main() {
	mylog.LogInputChan = mylog.LogStart()
	config.BuildConfig()
	global.Master.IP = config.Configs.Master_ip
	// ToDo:为客户端加入一张表，用来缓存用以沟通的数据表，其中相关的rpc连接要用的时候再去连
}
