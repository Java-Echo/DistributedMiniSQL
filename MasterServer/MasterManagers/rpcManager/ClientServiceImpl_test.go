package rpc

import (
	"testing"
	"time"
)

func TestCliService_Hello(t *testing.T) {
	go StartCliService("1234")
	// 延迟100秒，可别让它停下来了！
	time.Sleep(100 * time.Second)
}
