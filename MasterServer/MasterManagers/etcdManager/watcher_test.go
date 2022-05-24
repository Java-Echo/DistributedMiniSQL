package master

import (
	config "master/utils/ConfigSystem"
	"master/utils/global"
	"testing"
	"time"
)

func TestRegionRegisterWorker_OnPut(t *testing.T) {
	go RegisterWatcherWithWorker(global.Master, config.Configs.Etcd_region_register_catalog, &RegionRegisterWorker{})
	time.Sleep(30 * time.Second)
	t.Error("终止")
}
