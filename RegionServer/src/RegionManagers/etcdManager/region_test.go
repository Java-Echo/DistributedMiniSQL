package etcd

import (
	"fmt"
	config "region/utils/ConfigSystem"
	mylog "region/utils/LogSystem"
	"region/utils/global"
	"testing"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var cli *clientv3.Client

func TestMain(m *testing.M) {
	mylog.LogInputChan = mylog.LogStart()
	config.BuildConfig()
	global.Region = Init()
	global.MasterIP = GetMasterIP(global.Region)
	fmt.Println("初始化完成")
	m.Run()
}

// func TestInit(t *testing.T) {
// 	cli := Init()
// 	// go ServiceRegister(cli, "/server")
// 	value := GetMasterIP(cli)
// 	t.Log("master的IP为:" + value)
// 	go ServiceRegister(cli)
// 	t.Log("value=" + value)
// 	time.Sleep(1 * time.Second)
// 	t.Error("终止")
// }

func TestGetCopies(t *testing.T) {
	slaves := GetSlaves("bbb")
	for _, region := range slaves {
		fmt.Println(region)
	}
	t.Error("终止")
}

func TestGetSyncSlave(t *testing.T) {
	sync_slave := GetSyncSlave("ccc")
	fmt.Println(sync_slave)
	t.Error("终止")
}
