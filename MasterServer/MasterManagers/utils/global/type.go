package global

import (
	clientv3 "go.etcd.io/etcd/client/v3"
)

// region的信息表
type RegionMeta struct {
	IP    string
	Port  string
	State RegionState
}

type RegionState int32

const (
	Working RegionState = 0
	Stop    RegionState = 1
)

// table的信息表
type TableMeta struct {
	Name         string   // 表的名称
	MasterRegion string   // 主副本
	SyncRegion   string   // 同步从副本
	CopyRegions  []string // 异步从副本
}

// =========master的全局数据结构=========
var RegionMap map[string]*RegionMeta // 所有分区服务器的元信息
var TableMap map[string]*TableMeta   // 所有表的元信息
var HostIP string
var Master *clientv3.Client
