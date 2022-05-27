package regionWorker

import (
	"context"
	"fmt"
	"log"
	etcd "region/etcdManager"
	rpc "region/rpcManager"
	config "region/utils/ConfigSystem"
	mylog "region/utils/LogSystem"
	"region/utils/global"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 完成主从复制的相关逻辑
func CheckSlave(table global.TableMeta) (int, int) {
	// 首先清空当前的slave
	table.SyncRegion = ""
	table.CopyRegions = table.CopyRegions[0:0]
	// 连接etcd，得到该目录下的各个slave，将其填充到对应字段
	slaves := etcd.GetSlaves(table.Name)
	sync_slave := etcd.GetSyncSlave(table.Name)

	table.SyncRegion = sync_slave
	table.CopyRegions = slaves
	// 假设当前的slave不够，则会返回需求的slave个数
	res_sync := 0
	res_slave := 0
	if len(sync_slave) == 0 {
		res_sync = 1
	}
	if len(slaves) < 1 {
		res_slave = 1
	}
	return res_sync, res_slave
}

// ToDo:向master寻求slave
func GetSlave(tableName string, syncNeed int, slaveNeed int) {
	request := rpc.AskSlaveRst{
		TableName:    tableName,
		SyncSlaveNum: syncNeed,
		SlaveNum:     slaveNeed,
	}
	var reply rpc.AskSlaveRes
	rpc.RpcM2R.AskSlave(request, &reply)
	fmt.Println("GetSlave的返回值为:" + reply.State)
}

// ToDo:异步同步的日志管道处理
func StartAsyncCopy() chan<- global.SQLLog {
	// 这个是需要绑定在全局的
	input := make(chan global.SQLLog)
	go func() {
		for {
			log := <-input
			table := global.TableMap[log.Table]
			// 尝试获得写锁
			<-table.WriteLock
			fmt.Println("得到了表 '" + log.Table + "' 的写锁")
			rpc.MasterSQLChange(rpc.SQLRst{SQLtype: log.SQLtype, SQL: log.SQL, Table: log.Table})
			fmt.Println("成功执行了从副本同步：'" + log.SQL + "' ")
		}
	}()
	log_ := mylog.NewNormalLog("开启了全局异步SQL日志的复制")
	log_.LogType = "INFO"
	log_.LogGen(mylog.LogInputChan)
	return input
}

// ToDo:监听主副本在etcd上的目录的信息，根据这个目录的不同变化来做出不同的反应
func StartWatchTable(table *global.TableMeta) {
	catalog := config.Configs.Etcd_table_catalog + "/" + table.Name + "/"
	fmt.Println("监听的目录为:" + catalog)
	watchChan := global.Region.Watch(context.Background(), catalog, clientv3.WithPrefix())
	table.TableWatcher = &watchChan
	log_ := mylog.NewNormalLog("开启对表'" + table.Name + "'目录的监听")
	log_.LogGen(mylog.LogInputChan)

	for watchResponse := range watchChan {
		for _, event := range watchResponse.Events {
			ip := util_getLastKey(string(event.Kv.Key))
			level := util_getKey(string(event.Kv.Key), catalog, 0)

			if event.Type == clientv3.EventTypePut {
				// ToDo:此时有节点加入了，需要完成相应的逻辑
				fmt.Println("检测到表 '" + table.Name + "' 下有项目加入")
				fmt.Println("该表在本地为 '" + level + "' 类型的副本")
				fmt.Println("该表所对应的IP为 '" + ip + "' ")
				// 首先将其加入异步从副本，然后开启一个Goroutine向其传输日志文件快照(尽可能同时完成)
				go passTable(ip, table.Name)
				table.CopyRegions = append(table.CopyRegions, ip)
				log_ := mylog.NewNormalLog("成功为表 '" + table.Name + "' 添加了一个异步从副本 '" + ip + "' ")
				log_.LogGen(mylog.LogInputChan)
				// ToDo:如果是同步从副本的指令，则需要在日志全部运行完成的时候通知本程序，然后再将其加入到同步从副本中
			} else if event.Type == clientv3.EventTypeDelete {
				fmt.Println("检测到表 '" + table.Name + "' 下有项目删除")
				fmt.Println("该表在本地为 '" + util_getKey(string(event.Kv.Key), catalog, 0) + "' 类型的副本")
				fmt.Println("该表所对应的IP为 '" + util_getLastKey(string(event.Kv.Key)) + "' ")
				// 将节点从本地删去
				if level == "slave" {
					for i, ip_ := range table.CopyRegions {
						if ip_ == ip {
							table.CopyRegions = append(table.CopyRegions[:i], table.CopyRegions[i+1:]...)
							log_ := mylog.NewNormalLog("成功将表 '" + table.Name + "' 下的异步从副本 '" + ip + "' 删除")
							log_.LogGen(mylog.LogInputChan)
						}
					}
				} else if level == "sync_slave" {
					table.SyncRegion = ""
					log_ := mylog.NewNormalLog("成功将表 '" + table.Name + "' 下的同步从副本 '" + ip + "' 删除")
					log_.LogGen(mylog.LogInputChan)
				}
			}
		}
	}
}

// 同步从副本的监听机制，监听的是master目录的相关变化，一旦master目录变成了自己，那么自己就需要完成相关的提升
func StartWatchMaster(table *global.TableMeta) {
	catalog := config.Configs.Etcd_table_catalog + "/" + table.Name + "/master/"
	fmt.Println("监听的目录为:" + catalog)
	watchChan := global.Region.Watch(context.Background(), catalog, clientv3.WithPrefix())
	table.MasterWatcher = &watchChan
	log_ := mylog.NewNormalLog("开启对表'" + table.Name + "'的master目录的监听")
	log_.LogGen(mylog.LogInputChan)

	for watchResponse := range watchChan {
		for _, event := range watchResponse.Events {
			if event.Type == clientv3.EventTypePut {
				newMaster := util_getKey(string(event.Kv.Key), catalog, 0)
				fmt.Println("新成为master的节点为:" + newMaster)
				if newMaster == global.HostIP {
					fmt.Println("巧了,就是我要新成为master!")
					// 以下是在本地的表元信息映射表中，修改这张表
					table.Level = "master"
					catalog := config.Configs.Etcd_table_catalog + "/" + table.Name + "/"
					table.TableWatcher = etcd.GetWatcher(global.Region, catalog)
					go StartWatchTable(table)
					// 第一次遍历这个目录检查从副本数量，一旦数量不够，就向master索取
					syncNeed, slaveNeed := CheckSlave(*table)
					GetSlave(table.Name, syncNeed, slaveNeed)
				}
			} else if event.Type == clientv3.EventTypeDelete {
				oldMaster := util_getKey(string(event.Kv.Key), catalog, 0)
				fmt.Println("刚刚被删除的节点为:" + oldMaster)
				if oldMaster == global.HostIP {
					fmt.Println("糟了,被抛弃的节点竟然是我自己！")
					// ToDo:当前我们的设计是直接将被抛弃的主节点直接没了，连从副本都当不成
					// ToDo:尚未实现这里的逻辑
				}
			}
		}
	}
}

// ToDo:得到对应的表的内容
func getTableFile(tableName string) []byte {
	test := "123456"
	return []byte(test)
}

func passTable(dstIP string, tableName string) {
	fmt.Println("有被调用到！")
	client, err := rpc.DialGossipService("tcp", dstIP+":"+config.Configs.Rpc_R2R_port)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// 这里需要得到本地表格的byte数组
	var reply rpc.PassTableRes
	request := rpc.PassTableRst{
		Content:   getTableFile(tableName),
		TableName: tableName,
	}
	fmt.Println("已经做好了传输的准备")
	err_ := client.PassTable(request, &reply)
	fmt.Println("甚至成功返回了！")
	if err_ != nil {
		log.Fatal(err)
	} else {
		log_ := mylog.NewNormalLog("表 '" + tableName + "' 针对 '" + dstIP + "' 服务器的快照传输完成")
		log_.LogGen(mylog.LogInputChan)
	}
}
