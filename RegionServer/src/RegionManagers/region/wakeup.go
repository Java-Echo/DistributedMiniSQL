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

// 在本地的主副本文件夹下面查找所有的表名
func findMasterCopy(tableRoot string) []string {

	var files []string

	root := tableRoot + "/Master"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// 不处理文件夹
		if info.IsDir() {
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

// 在本地的从副本文件夹下面查找所有的表名
func findSlaveCopy(tableRoot string) []string {
	var files []string

	root := tableRoot + "/Slave"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// 不处理文件夹
		if info.IsDir() {
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
func SendNewTables(tableRoot string) {
	// ToDo:扫描本地的表的存储的文件夹得到本地的文件
	masterFiles := findMasterCopy(tableRoot)
	slaveFiles := findSlaveCopy(tableRoot)
	// ToDo:调用rpc接口进行发送
	var reply rpc.ReportTableRes
	var request []rpc.LocalTable
	request = make([]rpc.LocalTable, 0)
	for _, file := range masterFiles {
		table := rpc.LocalTable{}
		table.Name = file
		table.IP = global.HostIP
		table.Port = config.Configs.Rpc_R2R_port
		table.Level = "master"
		request = append(request, table)
	}
	for _, file := range slaveFiles {
		table := rpc.LocalTable{}
		table.Name = file
		table.IP = global.HostIP
		table.Port = config.Configs.Rpc_R2R_port
		table.Level = "slave"
		request = append(request, table)
	}
	err := rpc.RpcM2R.ReportTable(request, &reply)
	if err != nil {
		log.Fatal(err)
	}
	// ToDo:对rpc的返回值进行处理，建立本地的数据表
	// fmt.Println("本次返回的数组长度为:" + strconv.Itoa(len(reply.Tables)))
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
			go StartWatchTable(meta)
			// ToDo:第一次遍历这个目录检查从副本数量，一旦数量不够，就向master索取
			syncNeed, slaveNeed := CheckSlave(*meta)
			GetSlave(meta.Name, syncNeed, slaveNeed)
		} else if table.Level == "slave" {

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
