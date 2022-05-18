package master

import (
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	master := Init()
	go RegisterWatcher(master, "server/region_server")
	time.Sleep(10 * time.Second)
	t.Error("终止")
}
