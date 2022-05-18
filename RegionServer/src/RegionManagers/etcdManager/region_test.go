package etcd

import (
	"fmt"
	config "region/utils/ConfigSystem"
	mylog "region/utils/LogSystem"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var cli *clientv3.Client

func TestMain(m *testing.M) {
	cli = Init()
	mylog.LogInputChan = mylog.LogStart()
	config.BuildConfig()
	fmt.Println("初始化完成")
	m.Run()
}

func TestInit(t *testing.T) {
	cli := Init()
	// go ServiceRegister(cli, "/server")
	value := GetMasterAddress(cli)
	t.Log("master的IP为:" + value)
	go ServiceRegister(cli)
	t.Log("value=" + value)
	time.Sleep(1 * time.Second)
	t.Error("终止")
}
