package rpc

import (
	"log"
	mylog "master/utils/LogSystem"
	"net"
	"net/rpc"
)

type CliService struct{}

// 开启这个服务
func StartCliService(port string) {
	RegisterCliService(new(CliService))

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	log_ := mylog.NewNormalLog("开启了RPC(CliService)的监听服务,监听端口:" + port)
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
func (p *CliService) Hello(request string, reply *string) error {
	*reply = "你好"
	return nil
}
