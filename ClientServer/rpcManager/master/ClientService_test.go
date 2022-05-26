package masterRpc

import (
	config "client/utils/ConfigSystem"
	mylog "client/utils/LogSystem"
	"fmt"
	"log"
	"testing"
)

var client *CliServiceClient

func TestMain(m *testing.M) {
	mylog.LogInputChan = mylog.LogStart()
	config.BuildConfig()
	client, _ = DialService("tcp", "localhost:"+config.Configs.Master_port)
	m.Run()
}

func TestCliServiceClient_Hello(t *testing.T) {
	var request string
	var reply string
	err := client.Hello(request, &reply)
	if err != nil {
		log.Fatal(err)
	}
	t.Log("能够联系到master的服务")
	t.Error("终止")
}
func TestCliServiceClient_FetchTable(t *testing.T) {
	var request string
	request = "ttt"
	var reply TableInfo
	err := client.FetchTable(request, &reply)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("表名为:" + reply.Name)
	fmt.Println("这张表的主副本的IP地址为:" + reply.Master.IP)
	fmt.Println("这张表的同步从副本的IP地址为:" + reply.Sync_slave.IP)
	t.Error("终止")
}
