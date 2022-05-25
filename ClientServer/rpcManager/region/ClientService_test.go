package regionRpc

import (
	config "client/utils/ConfigSystem"
	mylog "client/utils/LogSystem"
	"log"
	"testing"
)

var client *CliServiceClient

func TestMain(m *testing.M) {
	mylog.LogInputChan = mylog.LogStart()
	config.BuildConfig()
	client, _ = DialService("tcp", "localhost:"+config.Configs.Region_port)
	m.Run()
}

func TestCliServiceClient_Hello(t *testing.T) {
	var request string
	var reply string
	err := client.Hello(request, &reply)
	if err != nil {
		log.Fatal(err)
	}
	t.Log("能够联系到region的服务")
	t.Error("终止")
}

func TestCliServiceClient_SQL(t *testing.T) {
	var request SQLRst
	request.SQL = "select * from test1"
	request.SQLtype = "insert"
	request.Table = "test1"
	var reply SQLRes
	err := client.SQL(request, &reply)
	if err != nil {
		log.Fatal(err)
	}
	t.Log("SQL语句的执行状态为:" + reply.State)
	t.Log("SQL语句的返回值为:" + reply.Result)
	t.Error("终止")
}
