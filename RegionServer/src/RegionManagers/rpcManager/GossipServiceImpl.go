package rpc

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
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
		SQLtype: request.SqlType,
		Table:   request.Table,
		SQL:     request.Sql,
	}
	log_ := mylog.NewNormalLog("接受到来自表" + request.Table + "的同步日志")
	log_.LogType = "INFO"
	log_.LogGen(mylog.LogInputChan)
	return nil
}

func (p *GossipService) PassTable(request PassTableRst, reply *PassTableRes) error {
	tableName := request.TableName

	// 1. 根据表的元信息进行相关的本地信息表的维护
	meta := &global.TableMeta{}
	meta.Name = tableName
	meta.Level = "slave"
	meta.WriteLock = make(chan int)

	// 2. 首先接受整个表文件
	file, err := os.Create(tableName + "_log11")
	if err != nil {
		log.Fatal("创建文件失败")
	}
	defer func() {
		file.Close()
	}()
	file.Write(request.Content)

	// 3. 尝试逐行读取其中的命令，然后执行SQL
	// br := bufio.NewReader(file)
	// for {
	// 	sqlLine, _, c := br.ReadLine()
	// 	if c == io.EOF {
	// 		break
	// 	}
	// 	fmt.Println(string(sqlLine))
	// 	// NormalSQL(string(sqlLine))
	// }
	// 4. 创建成功
	log_ := mylog.NewNormalLog("创建表 '" + tableName + "' 的备份成功")
	log_.LogType = "INFO"
	log_.LogGen(mylog.LogInputChan)

	// 3. 归还写锁(其实本来是没有的)
	meta.WriteLock <- 1
	fmt.Println("归还写锁")

	// 要不要开启一些channel之类的？
	return nil
}
