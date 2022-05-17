package rpc

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"strconv"
)

type ReportService struct{}

// ToDo:将得到的表的名字存储在本地的map当中
func (p *ReportService) ReportTable(request []LocalTable, reply *string) error {
	for _, t := range request {
		fmt.Print(t.Name + ":")
		fmt.Println(t.IP)
		// TableMap[t.Name] = t.IP
	}
	*reply = "你成功了,我一共接受到" + strconv.Itoa(len(request)) + "个表的信息"
	// fmt.Println("当前TableMap中一共有表项:" + strconv.Itoa(len(TableMap)) + "个")
	return nil
}

// 开启这个服务
func StartReportService() {
	RegisterReportService(new(ReportService))

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}

		go rpc.ServeConn(conn)
	}
}

// ToDo:完成更多的用于接受分区服务器报告的RPC调用
