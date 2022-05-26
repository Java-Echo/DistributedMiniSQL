package global

import (
	"fmt"

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

// =========测试函数=========

func PrintRegionMap(indent int) {
	if len(RegionMap) == 0 {
		fmt.Println("本地注册的region服务器为空！")
	}
	for _, region := range RegionMap {
		fmt.Println("------------------")
		PrintRegion(1, *region)
	}
}

func PrintTableMap(indent int) {
	if len(TableMap) == 0 {
		fmt.Println("本地的数据表为空！")
	}
	for _, table := range TableMap {
		fmt.Println("---------" + table.Name + "---------")
		PrintTable(1, *table)
	}
}

func PrintRegion(indent int, region RegionMeta) {
	fmt.Println(printIndent(indent) + "region ip:" + region.IP)
	var state string
	switch region.State {
	case Working:
		state = "working"
	case Stop:
		state = "stop"
	}
	fmt.Println(printIndent(indent) + "region state:" + state)
}

func PrintTable(indent int, table TableMeta) {
	fmt.Println(printIndent(indent) + "Name:" + table.Name)
	fmt.Println(printIndent(indent) + "Master:" + table.MasterRegion)
	fmt.Println(printIndent(indent) + "SyncRegion:" + table.SyncRegion)
	slaves := ""
	for _, ip := range table.CopyRegions {
		slaves += (ip + ",")
	}
	fmt.Println(printIndent(indent) + "SlaveRegions:" + slaves)
}

func printIndent(ind int) string {
	res := ""
	for i := 0; i < ind; i++ {
		res += "  "
	}
	return res
}
