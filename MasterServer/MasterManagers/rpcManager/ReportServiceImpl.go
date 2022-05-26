package rpc

import (
	"fmt"
	"log"
	etcd "master/etcdManager"
	mylog "master/utils/LogSystem"
	"master/utils/global"
	"net"
	"net/rpc"
)

type ReportService struct{}

// 开启这个服务
func StartReportService(port string) {
	RegisterReportService(new(ReportService))

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	log_ := mylog.NewNormalLog("开启了RPC(ReportService)的监听服务,监听端口:" + port)
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
		meta := &global.TableMeta{}
		meta.Name = t.Name
		// ToDo:暂时让当前服务器就作为该表的主副本节点，之后需要进一步的判断！
		// 在master本地更新相关信息
		meta.MasterRegion = t.IP
		global.TableMap[t.Name] = meta
		log := mylog.NewNormalLog("新增一张数据表:" + t.Name)
		log.LogType = "INFO"
		log.LogGen(mylog.LogInputChan)
		// 在etcd上创建该表的目录
		etcd.CreateMaster(meta.Name, t.IP)
		// 在返回值中指明“该表会成为主副本”
		table := ValidTable{}
		table.Level = "master"
		table.Name = t.Name
		reply.Tables = append(reply.Tables, table)
	}

	global.PrintTableMap(1)

	return nil
}

// ToDo:设立合适的选择region的逻辑，然后及时将这个region添加到etcd上面
// ToDo:尚未完成节点筛选逻辑的实现
func (p *ReportService) AskSlave(request AskSlaveRst, reply *AskSlaveRes) error {
	masterIP := global.TableMap[request.TableName].MasterRegion
	source := []string{masterIP}
	newSyncSlave := selectRegion(request.SyncSlaveNum, source)
	source = append(source, newSyncSlave...)
	newSlave := selectRegion(request.SlaveNum, source)
	// 在etcd上注册
	for _, s := range newSyncSlave {
		etcd.CreateSyncSlave(request.TableName, s)
		global.TableMap[request.TableName].SyncRegion = s
	}
	for _, s := range newSlave {
		etcd.CreateSlave(request.TableName, s)
		global.TableMap[request.TableName].CopyRegions = append(global.TableMap[request.TableName].CopyRegions, s)
	}

	reply.State = "成功！"
	global.PrintTableMap(1)
	return nil
}

// ToDo:从除了source之外的region中选出n个合适的
func selectRegion(n int, source []string) []string {
	// 检查一下source中的表
	fmt.Println("本地所有的表为:")
	for ip, _ := range global.RegionMap {
		fmt.Println(ip)
	}
	fmt.Println("source中的表为:")
	for _, ip := range source {
		fmt.Println(ip)
	}

	// 真正开始找表的地方
	res := make([]string, 0)
	num := 0
	for ip, _ := range global.RegionMap {
		if num >= n {
			break
		}
		inSource := false
		for _, ip_ := range source {
			if ip_ == ip {
				inSource = true
				break
			}
		}
		if !inSource {
			res = append(res, ip)
			num++
		}
	}
	fmt.Print("我们选定的表为:")
	fmt.Println(res)
	return res
}

// ToDo:完成更多的用于接受分区服务器报告的RPC调用
