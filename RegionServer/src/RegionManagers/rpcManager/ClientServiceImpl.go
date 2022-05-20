package rpc

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	mylog "region/utils/LogSystem"
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
		log_ = mylog.NewNormalLog("RPC服务(CliService)被调用")
		log_.LogType = "INFO"
		log_.LogGen(mylog.LogInputChan)
		go rpc.ServeConn(conn)
	}
}

//========具体的业务处理函数的实现========

// ToDo:更改处理逻辑
func (p *CliService) Hello(request string, reply *string) error {
	*reply = "你好"
	return nil
}

func (p *CliService) SQL(request SQLRst, reply *SQLRes) error {
	fmt.Println("接受到的SQL语句为:" + request.SQL)
	fmt.Println("SQL语句具体要作用的表为:" + request.Table)
	fmt.Println("接受到的SQL语句的类型为:" + request.SQLtype)
	return nil
}
