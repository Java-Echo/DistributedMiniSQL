package masterRpc

type Region struct {
	IP string
}

type TableInfo struct {
	Name       string
	Master     Region   // 主副本所在的节点
	Sync_slave Region   // 同步从副本所在的节点
	Slaves     []Region // 异步从副本所在的节点
}
