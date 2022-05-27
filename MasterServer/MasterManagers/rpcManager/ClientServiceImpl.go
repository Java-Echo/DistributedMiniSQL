package rpc

import (
	"fmt"
	"log"
	mylog "master/utils/LogSystem"
	"master/utils/global"
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
	table, ok := global.TableMap[request]
	if ok {
		// 此时本地能找到对应的数据表
		reply.Master.IP = table.MasterRegion
		reply.Sync_slave.IP = table.SyncRegion
		for _, ip := range table.CopyRegions {
			reply.Slaves = append(reply.Slaves, Region{IP: ip})
		}
		return nil
	}
	return fmt.Errorf("集群中未找到表 '" + request + "'")
}
