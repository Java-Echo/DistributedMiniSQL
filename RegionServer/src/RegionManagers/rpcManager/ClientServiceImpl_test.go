package rpc

import (
	"testing"
	"time"
)

func TestCliService_Hello(t *testing.T) {
	go StartCliService("5000")
	// 延迟100秒，可别让它停下来了！
	time.Sleep(30 * time.Second)
}
