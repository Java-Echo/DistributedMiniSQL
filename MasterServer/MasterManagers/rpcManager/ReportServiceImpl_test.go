package rpc

import (
	"master/utils/global"
	"testing"
	"time"
)

func TestStartReportService(t *testing.T) {
	global.TableMap = make(map[string]*global.TableMeta)
	go StartReportService("1237")
	time.Sleep(10 * time.Second)
	t.Error("终止")
}
