package rpc

import (
	"log"
	mylog "master/utils/LogSystem"
	"master/utils/global"
	"net"
	"net/rpc"
)

type ReportService struct{}

// 开启这个服务
func StartReportService() {
	RegisterReportService(new(ReportService))

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	log_ := mylog.NewNormalLog("开启了RPC(ReportService)的监听服务")
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

// ToDo:将得到的表的名字存储在本地的map当中
func (p *ReportService) ReportTable(request []LocalTable, reply *ReportTableRes) error {
	// 检查：是否已经存在同名表
	for _, t := range request {
		if _, ok := global.TableMap[t.Name]; ok {
			// 此时检测到同名表
			// ToDo:这里得返回错误信息，错误系统我尚未建立
			return nil
		}
	}

	// 将分区服务器的所有表加入本地，同时令该服务器为表的master
	for _, t := range request {
		meta := global.TableMeta{}
		meta.Name = t.Name
		// ToDo:暂时让当前服务器就作为该表的主副本节点，之后需要进一步的判断！
		// 在master本地更新相关信息
		meta.MasterRegion = t.IP
		global.TableMap[t.Name] = meta
		log := mylog.NewNormalLog("新增一张数据表:" + t.Name)
		log.LogType = "INFO"
		log.LogGen(mylog.LogInputChan)
		// 在返回值中指明“该表会成为主副本”
		table := ValidTable{}
		table.Level = "master"
		table.Name = t.Name
		reply.Tables = append(reply.Tables, table)
	}
	return nil
}

// ToDo:完成更多的用于接受分区服务器报告的RPC调用
