package global

//==========client的全局数据结构==========
type Region struct {
	IP string
	// 当前的设计中，分区服务器对client暴露的端口是固定的，所以暂时不需要存储端口
}
type TableMeta struct {
	Name       string
	Master     Region   // 主副本所在的节点
	Sync_slave Region   // 同步从副本所在的节点
	Slaves     []Region // 异步从副本所在的节点
}

//==========client的全局数据==========
var TableCache []TableMeta
var Master Region
