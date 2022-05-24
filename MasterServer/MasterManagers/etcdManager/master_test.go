package master

import (
	"fmt"
	config "master/utils/ConfigSystem"
	mylog "master/utils/LogSystem"
	"master/utils/global"
	"testing"
)

func TestMain(m *testing.M) {
	mylog.LogInputChan = mylog.LogStart()
	config.BuildConfig()
	global.Master = Init()
	global.RegionMap = make(map[string]*global.RegionMeta)
	global.TableMap = make(map[string]*global.TableMeta)
	fmt.Println("初始化完成")
	m.Run()
}

// func TestRegisterWatcher(t *testing.T) {
// 	// 开启日志功能
// 	go RegisterWatcher(master, "/server/")
// 	time.Sleep(10 * time.Second)
// 	t.Error("终止")
// }

func TestCreateSlave(t *testing.T) {
	CreateSlave("bbb", "123.456.789.0")

	t.Error("终止")
}

func TestCreateSyncCopys(t *testing.T) {
	CreateSyncSlave("ccc", "987,654,321,0")
	t.Error("终止")
}
