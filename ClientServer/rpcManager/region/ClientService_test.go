package regionRpc

import (
	config "client/utils/ConfigSystem"
	mylog "client/utils/LogSystem"
	"log"
	"testing"
)

func TestMain(m *testing.M) {
	mylog.LogInputChan = mylog.LogStart()
	config.BuildConfig()
	m.Run()
}

func TestCliServiceClient_Hello(t *testing.T) {
	client, err := DialService("tcp", "localhost:"+config.Configs.Region_port)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	var request string
	var reply string
	err = client.Hello(request, &reply)
	if err != nil {
		log.Fatal(err)
	}
	t.Log("能够联系到region的服务")
	t.Error("终止")
}
