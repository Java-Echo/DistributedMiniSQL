package rpc

import (
	"testing"
	"time"
)

func TestStartReportService(t *testing.T) {
	StartReportService()
	time.Sleep(10 * time.Second)
	t.Error("终止")
}
