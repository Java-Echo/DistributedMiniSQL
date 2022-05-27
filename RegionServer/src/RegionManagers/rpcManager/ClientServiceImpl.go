package rpc

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	config "region/utils/ConfigSystem"
	mylog "region/utils/LogSystem"
	"region/utils/global"
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
	// 1. 首先检查本地是否有这张表
	table, ok := global.TableMap[request.Table]
	if !ok {
		reply.State = "error"
		reply.Result = "未在" + global.HostIP + "找到表 '" + request.Table + "' "
		log_ := mylog.NewNormalLog("未在 '" + global.HostIP + "' 找到表 '" + request.Table + "' ")
		log_.LogType = "ERROR"
		log_.LogGen(mylog.LogInputChan)
		return nil
	}
	// 2. 尝试执行其中的指令
	switch request.SQLtype {
	case "select":
		// 检查这张表的版本是否有问题(暂时先不做)
		// 调用sql的查询
		res, success := MasterSQLSelect(request)
		reply.Result = res
		if success {
			reply.State = "success"
			log_ := mylog.NewNormalLog("执行SQL语句 '" + request.SQL + "' 成功")
			log_.LogType = "INFO"
			log_.LogGen(mylog.LogInputChan)
			return nil
		} else {
			reply.State = "fail"
			log_ := mylog.NewNormalLog("执行SQL语句 '" + request.SQL + "' 失败")
			log_.LogType = "ERROR"
			log_.LogGen(mylog.LogInputChan)
			return nil
		}
	case "delete", "insert", "update":
		// 1. 首先检查本地是否有这张表，并查看该表的副本等级
		if table.Level != "master" {
			reply.State = "error"
			reply.Result = "分区服务器 '" + global.HostIP + "' 并不是表 '" + request.Table + "' 的主副本"
			return nil
		}
		// 2. 尝试在本地完成修改
		_, ok := MasterSQLChange(request)
		if ok {
			log_ := mylog.NewNormalLog("执行SQL语句 '" + request.SQL + "' 成功")
			log_.LogType = "INFO"
			log_.LogGen(mylog.LogInputChan)
		} else {
			log_ := mylog.NewNormalLog("执行SQL语句 '" + request.SQL + "' 失败")
			log_.LogType = "ERROR"
			log_.LogGen(mylog.LogInputChan)
		}
		// ToDo:尝试向同步从副本进行修改
		if table.SyncRegion != "" {
			// 执行同步修改操作
			fmt.Println("这里需要对 '" + table.SyncRegion + "' 进行同步修改")
		}
		// ToDo:尝试将相关信息存储到异步从副本当中
		for _, ip := range table.CopyRegions {
			fmt.Println("这里需要对 '" + ip + "' 进行异步修改")
			client, err := DialGossipService("tcp", ip+":"+config.Configs.Rpc_R2R_port)
			if err != nil {
				log.Fatal("dialing:", err)
			}
			var reply PassLogRes
			err = client.PassLog(PassLogRst{SqlType: request.SQLtype, Sql: request.SQL, Table: request.Table}, &reply)
			if err != nil {
				log.Fatal(err)
			}
		}
		// 添加返回值信息
		reply.State = "success"
		reply.Result = ""
		// 成功返回
		fmt.Println("对数据表的局部改动操作完成")
		return nil
	case "create_table":
		// ToDo:尝试向master申请主副本

		// 如果master成功返回，则在在本地进行SQL的执行，得到执行结果
		res, success := MasterSQLTableCreate(request)
		reply.Result = res
		if success {
			reply.State = "success"
			log_ := mylog.NewNormalLog("执行SQL语句 '" + request.SQL + "' 成功")
			log_.LogType = "INFO"
			log_.LogGen(mylog.LogInputChan)
			return nil
		} else {
			reply.State = "fail"
			log_ := mylog.NewNormalLog("执行SQL语句 '" + request.SQL + "' 失败")
			log_.LogType = "ERROR"
			log_.LogGen(mylog.LogInputChan)
			return nil
		}
	case "delete_table":
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
