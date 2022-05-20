package rpc

import (
	"log"
	config "region/utils/ConfigSystem"
	"testing"
)

func TestGossipService_PassLog(t *testing.T) {
	client, err := DialGossipService("tcp", "localhost:"+config.Configs.Rpc_R2R_port)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	var reply PassLogRes
	request := PassLogRst{
		SQLtype: "insert",
		SQL:     "inset xxx",
		Table:   "aaa",
	}
	err = client.PassLog(request, &reply)
	if err != nil {
		log.Fatal(err)
	}
	t.Log("成功执行了应该，我猜的")
	t.Error("终止")
}
