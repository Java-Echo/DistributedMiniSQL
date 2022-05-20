package regionWorker

import (
	"fmt"
	etcd "region/etcdManager"
	rpc "region/rpcManager"
	mylog "region/utils/LogSystem"
	"region/utils/global"
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
			fmt.Println("要同步的log为:" + log.SQL)
		}
	}()
	log_ := mylog.NewNormalLog("开启了全局异步SQL日志的复制")
	log_.LogType = "INFO"
	log_.LogGen(mylog.LogInputChan)
	return input
}
