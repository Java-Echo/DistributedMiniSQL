package regionWorker

import (
	etcd "region/etcdManager"
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

// 向master寻求slave
func GetSlave(syncNeed int, slaveNeed int) {

}
