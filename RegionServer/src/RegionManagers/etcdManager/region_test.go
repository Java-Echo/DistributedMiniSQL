package etcd

import (
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	cli := Init()
	// go ServiceRegister(cli, "/server")
	value := GetMasterConfig(cli, "/config/masterAddress")
	t.Log("value=" + value)
	time.Sleep(1 * time.Second)
	t.Error("终止")
}
