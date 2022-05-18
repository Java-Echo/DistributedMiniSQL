package rpc

import (
	"fmt"
	"log"
	config "region/utils/ConfigSystem"
	mylog "region/utils/LogSystem"
	"testing"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var cli *clientv3.Client

func TestMain(m *testing.M) {
	mylog.LogInputChan = mylog.LogStart()
	config.BuildConfig()
	fmt.Println("初始化完成")
	m.Run()
}

func TestReportServiceClient_ReportTable(t *testing.T) {
	client, err := DialReportService("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply string
	request := []LocalTable{
		{"aab", "127.0.0.1", "1234"},
		{"bbc", "127.0.0.1", "1234"},
		{"cca", "127.0.0.1", "1234"},
	}
	err = client.ReportTable(request, &reply)
	if err != nil {
		log.Fatal(err)
	}
	t.Log(reply)
	t.Error("终止")
}
