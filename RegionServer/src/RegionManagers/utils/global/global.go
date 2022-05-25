package global

import (
	clientv3 "go.etcd.io/etcd/client/v3"
)

type SQLLog struct {
	SQLtype string // SQL语句的类型
	Table   string // SQL语句具体查询的表
	SQL     string // 具体的SQL语句
}
type TableMeta struct {
	Name         string              // 表的名称
	Level        string              // 表的等级(master/slave/sync_slave)
	State        string              // 表的状态
	WriteLock    chan int            // 写锁
	TableWatcher *clientv3.WatchChan // 监听表在etcd上的目录(只有在等级为master的时候有用)
	SyncRegion   string              // 同步从副本(只有在等级为master的时候有用)
	CopyRegions  []string            // 异步从副本(只有在等级为master的时候有用)
}

//==========region的全局数据结构==========
var Region *clientv3.Client        // 连接etcd的节点
var TableMap map[string]*TableMeta // 所有表的元信息
var MasterIP string                // master节点的IP地址
var HostIP string                  // 本地的IP地址
var AsyncLogSQLChan chan<- SQLLog  //全局的SQL异步备份管道
