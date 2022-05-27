package main

import (
	masterRPC "client/rpcManager/master"
	regionRPC "client/rpcManager/region"
	config "client/utils/ConfigSystem"
	mylog "client/utils/LogSystem"
	"client/utils/global"
	"fmt"
	"log"
	"strings"
)

func main() {
	mylog.LogInputChan = mylog.LogStart()
	config.BuildConfig()
	global.Master.IP = config.Configs.Master_ip
	// 注册master端的rpc服务
	masterRPC.RpcM2R, _ = masterRPC.DialService("tcp", config.Configs.Master_ip+":"+config.Configs.Master_port)
	// ToDo:为客户端加入一张表，用来缓存用以沟通的数据表，其中相关的rpc连接要用的时候再去连

	for {
		var sql string
		fmt.Print("sql>>>")
		fmt.Scanln(&sql)
		res, _ := runSQL(sql)
		fmt.Println(res)
	}

}

// ToDo:直接返回一个查询体，主要需要解析出来 ①查询的table名称 ②执行的操作类型
func parser(input string) (regionRPC.SQLRst, bool) {
	rst := regionRPC.SQLRst{}
	// [检查]是否以分号结尾
	if input[len(input)-1] != ';' {
		fmt.Println("SQL语句没有以分号结尾!")
		return rst, false
	}
	// 将分号删除
	input = input[:len(input)-1]
	// 分割字段
	word := strings.Split(input, " ")
	input += ";"
	// 通过第一个word判断操作类型
	if word[0] == "select" || word[0] == "insert" || word[0] == "delete" || word[0] == "update" {
		rst.SQLtype = word[0]
		rst.SQL = input
		// ToDo:结合特定的语句获得table的名称
		switch rst.SQLtype {
		case "select":
			// select column1 from table1 where
			for i, s := range word {
				if s == "from" {
					rst.Table = word[i+1]
					break
				}
			}
		case "insert":
			// insert into table1 values
			rst.Table = word[2]
		case "delete":
			// delete from student2 where id=1080100245;
			rst.Table = word[2]
		case "update":
			// update student2 set
			rst.Table = word[1]
		}
	} else if word[0] == "drop" {
		//  假定这两个的table字段都是在第三位
		rst.SQL = input
		rst.Table = word[2]
		rst.SQLtype = "drop_table"
	} else if word[0] == "create" {
		//  假定这两个的table字段都是在第三位
		rst.SQL = input
		rst.Table = word[2]
		rst.SQLtype = "create_table"
	} else {
		fmt.Println("错误的操作符!")
		return rst, false
	}
	return rst, true
}

// 真正尝试运行SQL的程序
func runSQL(input string) (string, error) {
	// 1. 得到解析后的SQL内容
	sqlRst, ok := parser(input)
	fmt.Println("SQL解析完成")
	if ok {
		// 2. 此时SQL不会有问题，尝试在本地查找相关的表
		var table global.TableMeta // 存储得到的表的信息
		var inCache bool
		table, inCache = global.TableCache[sqlRst.Table]
		if !inCache {
			fmt.Println("表 '" + sqlRst.Table + "' 不在缓存中,我们需要向master询问")
			// 此时table不在缓存中，我们需要向master询问
			reply := masterRPC.TableInfo{}
			err := masterRPC.RpcM2R.FetchTable(sqlRst.Table, &reply)
			if err != nil {
				fmt.Println("未在集群中找到表 '" + sqlRst.Table + "' ")
				return "", fmt.Errorf("未在集群中找到表 '" + sqlRst.Table + "' ")
			}
			table.Master.IP = reply.Master.IP
			fmt.Println("找到的MasterIP:" + reply.Master.IP)
			table.Sync_slave.IP = reply.Sync_slave.IP
			fmt.Println("找到的Sync_slaveIP:" + reply.Sync_slave.IP)
			for _, slave := range reply.Slaves {
				table.Slaves = append(table.Slaves, global.RegionInfo{IP: slave.IP})
				fmt.Println("找到的SlavesIP:" + slave.IP)
			}
			// 将该表加入缓存
			global.TableCache[sqlRst.Table] = table
		}
		// 3. 此时理论上我们已经获得了存储有这张表对应的region服务器
		table = global.TableCache[sqlRst.Table] // 应该是要重新赋值一下的
		res, err := chooseRegionAndRun(sqlRst, table)
		if err != nil {
			log.Fatal("runSQL error:", err)
			return "", err
		} else {
			// 此时终于没有问题了
			return res, nil
		}

	} else {
		fmt.Println("错误的SQL语句")
		return "", fmt.Errorf("错误的SQL语句")
	}
}

// 尝试在当前保存的表的region信息中选择一个region服务器进行运行
func chooseRegionAndRun(sql regionRPC.SQLRst, tableMeta global.TableMeta) (string, error) {
	printTableMeta(tableMeta)
	// 尝试一个个进行tableMeta进行尝试连接
	if sql.SQLtype == "select" {
		// 优先级①:尝试向slaves获取内容
		for _, slave := range tableMeta.Slaves {
			if slave.IP != "" {
				fmt.Println("将尝试向服务器 '" + slave.IP + "' 发送访问请求")
				res, err := runOnRegion(sql, slave.IP)
				if err == nil {
					log_ := mylog.NewNormalLog("sql语句 '" + sql.SQL + "' 将在服务器 '" + slave.IP + "' 运行")
					log_.LogType = "INFO"
					log_.LogGen(mylog.LogInputChan)
					return res, nil
				}
			}
		}
		// 优先级②:尝试向sync_slave获取内容
		if tableMeta.Sync_slave.IP != "" {
			fmt.Println("将尝试向服务器 '" + tableMeta.Sync_slave.IP + "' 发送访问请求")
			res, err := runOnRegion(sql, tableMeta.Sync_slave.IP)
			if err == nil {
				log_ := mylog.NewNormalLog("sql语句 '" + sql.SQL + "' 将在服务器 '" + tableMeta.Sync_slave.IP + "' 运行")
				log_.LogType = "INFO"
				log_.LogGen(mylog.LogInputChan)
				return res, nil
			}
		}
		// 优先级③:尝试向master获取内容
		if tableMeta.Master.IP != "" {
			fmt.Println("将尝试向服务器 '" + tableMeta.Master.IP + "' 发送访问请求")
			res, err := runOnRegion(sql, tableMeta.Master.IP)
			if err == nil {
				log_ := mylog.NewNormalLog("sql语句 '" + sql.SQL + "' 将在服务器 '" + tableMeta.Master.IP + "' 运行")
				log_.LogType = "INFO"
				log_.LogGen(mylog.LogInputChan)
				return res, nil
			}
		}
	} else {
		// 此时只能够向master传递请求
		if tableMeta.Master.IP != "" {
			fmt.Println("将尝试向服务器 '" + tableMeta.Master.IP + "' 发送访问请求")
			res, err := runOnRegion(sql, tableMeta.Master.IP)
			if err == nil {
				log_ := mylog.NewNormalLog("sql语句 '" + sql.SQL + "' 将在服务器 '" + tableMeta.Master.IP + "' 运行")
				log_.LogType = "INFO"
				log_.LogGen(mylog.LogInputChan)
				return res, nil
			} else {
				return "", err
			}
		}
	}
	return "", fmt.Errorf("当前缓冲中的所有region都无法访问")
}

// 尝试堆特定的IP地址进行特定的访问(已测试)
func runOnRegion(sql regionRPC.SQLRst, clientIP string) (string, error) {
	// 尝试联络
	fmt.Println("尝试连接服务器:" + clientIP)
	client, err := regionRPC.DialService("tcp", clientIP+":"+config.Configs.Region_port)
	if err != nil {
		return "", err
	}
	fmt.Println("成功与 '" + clientIP + "' 的主机建立联系,尝试从它的地方获取数据")
	// 尝试传输sql的请求
	var reply regionRPC.SQLRes
	err = client.SQL(sql, &reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("成功的")

	// 接受相关的请求
	fmt.Println("SQL语句的执行状态为:" + reply.State)
	fmt.Println("SQL语句的返回值为:" + reply.Result)
	if reply.State == "fail" {
		return reply.Result, fmt.Errorf("SQL执行失败")
	} else {
		return reply.Result, nil
	}
}

func printTableMeta(table global.TableMeta) {
	fmt.Println("-------" + table.Name + "-------")
	fmt.Println("table master:" + table.Master.IP)
	fmt.Println("table sync_slave:" + table.Sync_slave.IP)
	slaves := ""
	for _, slave := range table.Slaves {
		slaves += slave.IP + ","
	}
	fmt.Println("table slaves:" + slaves)
	fmt.Println("-------------------")
}

func printSQL(sql regionRPC.SQLRst) {
	fmt.Println("---------sql---------")
	fmt.Println("sql:" + sql.SQL)
}
