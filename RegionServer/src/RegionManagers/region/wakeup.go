package regionWorker

import (
	"log"
	"os"
	"path/filepath"
	etcd "region/etcdManager"
	rpc "region/rpcManager"
	config "region/utils/ConfigSystem"
	"region/utils/global"
	"strings"
)

// // 在本地的主副本文件夹下面查找所有的表名
// func findMasterCopy(tableRoot string) []string {

// 	var files []string

// 	root := tableRoot + "/Master"
// 	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
// 		// 不处理文件夹
// 		if info.IsDir() {
// 			return nil
// 		}
// 		// 将搜索到的文件名添加进来
// 		files = append(files, info.Name())
// 		return nil
// 	})
// 	if err != nil {
// 		panic(err)
// 	}

// 	return files
// }

// // 在本地的从副本文件夹下面查找所有的表名
// func findSlaveCopy(tableRoot string) []string {
// 	var files []string

// 	root := tableRoot + "/Slave"
// 	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
// 		// 不处理文件夹
// 		if info.IsDir() {
// 			return nil
// 		}
// 		// 将搜索到的文件名添加进来
// 		files = append(files, info.Name())
// 		return nil
// 	})
// 	if err != nil {
// 		panic(err)
// 	}

// 	return files
// }

func findLocalTable(tableRoot string) []string {
	var files []string

	root := tableRoot
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// 不处理文件夹
		if info.IsDir() {
			return nil
		}
		fileName := info.Name()
		prefix := "_list"
		if len(fileName) >= len(prefix) && fileName[len(fileName)-len(prefix):] == prefix {
			return nil
		}
		// 将搜索到的文件名添加进来
		files = append(files, info.Name())
		return nil
	})
	if err != nil {
		panic(err)
	}

	return files
}

// 完成分区服务器新建和重启的相关任务
func SendLocalTables(tableRoot string) {
	// 1.扫描本地的表的存储的文件夹得到本地的文件
	masterFiles := findLocalTable(tableRoot)
	// 2.
	tables := make([]rpc.LocalTable, 0)
	for _, file := range masterFiles {
		table := rpc.LocalTable{}
		table.Name = file
		table.IP = global.HostIP
		table.Port = config.Configs.Rpc_R2R_port
		table.Level = "master"
		tables = append(tables, table)
	}
	SendNewTables(tables)

	global.PrintTableMap(1)
}

func SendNewTables(tables []rpc.LocalTable) {
	var reply rpc.ReportTableRes
	err := rpc.RpcM2R.ReportTable(tables, &reply)
	if err != nil {
		log.Fatal(err)
	}
	for _, table := range reply.Tables {
		meta := &global.TableMeta{}
		meta.Level = table.Level
		meta.Name = table.Name
		if table.Level == "master" {
			// 建立这个表的其他元信息
			catalog := config.Configs.Etcd_table_catalog + "/" + table.Name + "/"
			meta.TableWatcher = etcd.GetWatcher(global.Region, catalog)
			global.TableMap[meta.Name] = meta
			// ToDo:完善对于主副本建立监听机制
			go rpc.StartWatchTable(meta)
			// ToDo:第一次遍历这个目录检查从副本数量，一旦数量不够，就向master索取
			syncNeed, slaveNeed := rpc.CheckSlave(*meta)
			rpc.GetSlave(meta.Name, syncNeed, slaveNeed)
		}
	}
}

// 工具函数：得到路径的最后一个字段
func util_getLastKey(path string) string {
	keys := strings.Split(path, "/")
	return keys[len(keys)-1]
}

func util_getKey(path string, prefix string, n int) string {
	str := path[len(prefix):]
	keys := strings.Split(str, "/")
	return keys[n]
}
