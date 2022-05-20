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

func (p *CliService) FetchTable(request string, reply *TableInfo) error {
	// ToDo:尝试查表，并将相关的信息添加到TableInfo中
	reply.Master = Region{"master的IP地址"}
	reply.Sync_slave = Region{"sync_slave的IP地址"}
	reply.Slaves = append(reply.Slaves, Region{"第一个异步从节点的IP地址"})
	reply.Slaves = append(reply.Slaves, Region{"第二个异步从节点的IP地址"})
	return nil
}
