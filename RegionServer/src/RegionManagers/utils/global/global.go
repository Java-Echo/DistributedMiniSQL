package global

import (
	rpc "region/rpcManager"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type TableMeta struct {
	Name         string              // 表的名称
	Level        string              // 表的等级(master/slave/sync_slave)
	State        string              // 表的状态
	TableWatcher *clientv3.WatchChan // 监听表在etcd上的目录(只有在等级为master的时候有用)
	SyncRegion   string              // 同步从副本(只有在等级为master的时候有用)
	CopyRegions  [10]string          // 异步从副本(只有在等级为master的时候有用)
}

//==========region的全局数据结构==========
var Region *clientv3.Client         // 连接etcd的节点
var RpcM2R *rpc.ReportServiceClient // 与master通信的rpc服务
var TableMap map[string]TableMeta   // 所有表的元信息
var MasterIP string                 // master节点的IP地址
var HostIP string                   // 本地的IP地址
