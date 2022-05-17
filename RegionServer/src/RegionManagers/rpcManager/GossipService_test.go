package rpc

import (
	"log"
	"testing"
	"time"
)

func TestGossipServiceClient_FetchLog(t *testing.T) {
	go StartGossipService("1234")
	time.Sleep(1 * time.Second)

	client, err := DialGossipService("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// 测试第一个方法
	var reply FetchLogRes
	request := FetchLogRst{}
	request.TableName = "table_test"
	request.Version = "version_1"
	err = client.FetchLog(request, &reply)
	if err != nil {
		log.Fatal(err)
	}
	t.Log(reply)
	time.Sleep(1 * time.Second)

	// 测试第二个方法
	var reply2 SyncProbeRes
	request2 := SyncProbeRst{}
	request2.TableName = "table_test2"
	request2.Version = "version_2"
	err = client.SyncProbe(request2, &reply2)
	if err != nil {
		log.Fatal(err)
	}
	t.Log(reply2)
	time.Sleep(1 * time.Second)

	t.Error("终止")
}
