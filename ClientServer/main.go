package main

import (
	config "client/utils/ConfigSystem"
	mylog "client/utils/LogSystem"
	"client/utils/global"
	"strings"
)

func main() {
	mylog.LogInputChan = mylog.LogStart()
	config.BuildConfig()
	global.Master.IP = config.Configs.Master_ip
	// ToDo:为客户端加入一张表，用来缓存用以沟通的数据表，其中相关的rpc连接要用的时候再去连
}

// ToDo:直接返回一个查询体
func parser(input string) {
	word := strings.Split(input, " ")
	switch word[0] {
	case "select":
	case "insert":
	case "delete":

	}
}
