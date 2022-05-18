package master

import (
	"fmt"
	mylog "master/utils/LogSystem"
	global "master/utils/global"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var master *clientv3.Client

func TestMain(m *testing.M) {
	master = Init()
	mylog.LogInputChan = mylog.LogStart()
	global.RegionMap = make(map[string]global.RegionMeta)
	m.Run()
	fmt.Println("初始化完成")
}

func TestRegisterWatcher(t *testing.T) {
	// 开启日志功能
	go RegisterWatcher(master, "/server/")
	time.Sleep(10 * time.Second)
	t.Error("终止")
}
