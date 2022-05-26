package global

import (
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type SQLLog struct {
	SQLtype string // SQL语句的类型
	Table   string // SQL语句具体查询的表
	SQL     string // 具体的SQL语句
}
type TableMeta struct {
	Name          string              // 表的名称
	Level         string              // 表的等级(master/slave/sync_slave)
	State         string              // 表的状态
	WriteLock     chan int            // 写锁
	TableWatcher  *clientv3.WatchChan // 监听表在etcd上的目录(只有在等级为master的时候有用)
	MasterWatcher *clientv3.WatchChan // 监听表在etcd上的目录(只有等级为slave/sync_slave的有)
	SyncRegion    string              // 同步从副本(只有在等级为master的时候有用)
	CopyRegions   []string            // 异步从副本(只有在等级为master的时候有用)
}

//==========region的全局数据结构==========
var Region *clientv3.Client        // 连接etcd的节点
var TableMap map[string]*TableMeta // 所有表的元信息
var MasterIP string                // master节点的IP地址
var HostIP string                  // 本地的IP地址
var AsyncLogSQLChan chan<- SQLLog  //全局的SQL异步备份管道
var SQLInput chan string
var SQLOutput chan string

//==========工具函数==========
func PrintTableMap(indent int) {
	for _, table := range TableMap {
		fmt.Println("---------" + table.Name + "---------")
		PrintTableMeta(1, *table)
	}
}

func PrintTableMeta(indent int, table TableMeta) {
	fmt.Println(printIndent(indent) + "Name:" + table.Name)
	fmt.Println(printIndent(indent) + "Level:" + table.Level)
	fmt.Println(printIndent(indent) + "State:" + table.State)
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
