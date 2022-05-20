package rpc

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	mylog "region/utils/LogSystem"
)

type GossipService struct{}

// 开启这个服务
func StartGossipService(port string) {
	RegisterGossipService(new(GossipService))

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	log_ := mylog.NewNormalLog("开启了RPC(GossipService)的监听服务,监听端口:" + port)
	log_.LogType = "INFO"
	log_.LogGen(mylog.LogInputChan)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}

		go rpc.ServeConn(conn)
	}
}

//========具体的业务处理函数的实现========

// ToDo:更改处理逻辑
func (p *GossipService) FetchLog(request FetchLogRst, reply *FetchLogRes) error {
	fmt.Println("你的表名为" + request.TableName)
	fmt.Println("你的版本号为" + request.Version)
	reply.Log = []string{
		"111",
		"222",
		"333",
	}
	reply.TableName = "表的名字"
	reply.Version = "版本号"
	return nil
}

// ToDo:更改处理逻辑
func (p *GossipService) SyncProbe(request SyncProbeRst, reply *SyncProbeRes) error {
	fmt.Println("你的表名啊,它是" + request.TableName)
	fmt.Println("你想知道版本号" + request.Version + "我更新了没")

	return nil
}
