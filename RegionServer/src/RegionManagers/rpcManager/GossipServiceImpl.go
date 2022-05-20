package rpc

import (
	"log"
	"net"
	"net/rpc"
	mylog "region/utils/LogSystem"
	"region/utils/global"
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
func (p *GossipService) PassLog(request PassLogRst, reply *PassLogRes) error {
	// 需要接受当前的日志
	global.AsyncLogSQLChan <- global.SQLLog{
		SQLtype: request.SQLtype,
		Table:   request.Table,
		SQL:     request.SQL,
	}
	return nil
}
