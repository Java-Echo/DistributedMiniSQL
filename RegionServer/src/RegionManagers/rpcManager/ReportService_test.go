package rpc

import (
	"fmt"
	"log"
	config "region/utils/ConfigSystem"
	mylog "region/utils/LogSystem"
	"strconv"
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
	client, err := DialReportService("tcp", "localhost:"+config.Configs.Rpc_M2R_port)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply ReportTableRes
	request := []LocalTable{
		{"aab", "127.0.0.1", "1234", ""},
		{"bbc", "127.0.0.1", "1234", ""},
		{"cca", "127.0.0.1", "1234", ""},
	}
	err = client.ReportTable(request, &reply)
	if err != nil {
		log.Fatal(err)
	}
	t.Log("返回的数组长度为:" + strconv.Itoa(len(reply.Tables)))
	t.Error("终止")
}

func TestReportServiceClient_AskSlave(t *testing.T) {
	client, err := DialReportService("tcp", "localhost:"+config.Configs.Rpc_M2R_port)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var request AskSlaveRst
	var reply AskSlaveRes
	request.TableName = "aaa"
	request.SyncSlaveNum = 1
	request.SlaveNum = 1
	err = client.AskSlave(request, &reply)
	if err != nil {
		log.Fatal(err)
	}
	t.Log("请求的状态为:" + reply.State)
	t.Error("终止")
}
