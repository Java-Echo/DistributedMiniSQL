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
	switch request.SQLtype {
	case "select":
		// 首先检查本地是否有这张表

		// 检查这张表的版本是否有问题(暂时先不做)

		// 调用sql的查询
		fmt.Println("单纯的查询操作")
	case "delete", "insert", "update":
		// 首先检查本地是否有这张表，并查看该表的副本等级

		// 尝试在本地完成修改

		// 尝试向同步从副本进行修改

		// 尝试将相关信息存储到异步从副本当中

		// 成功返回
		fmt.Println("对数据表的局部改动操作")
	case "create_table", "delete_table":
		// 首先在本地进行SQL的执行，得到执行结果

		// 本地创建成功，尝试向master申请主副本

		fmt.Println("这个操作不得了，要对数据表整体改动")
	default:
		fmt.Println("什么b操作?")
	}
	reply.Result = "什么都没有查到哦"
	reply.State = "成功"
	return nil
}
