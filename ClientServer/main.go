package main

import (
	regionRPC "client/rpcManager/region"
	config "client/utils/ConfigSystem"
	mylog "client/utils/LogSystem"
	"client/utils/global"
	"fmt"
	"strings"
)

func main() {
	mylog.LogInputChan = mylog.LogStart()
	config.BuildConfig()
	global.Master.IP = config.Configs.Master_ip
	// ToDo:为客户端加入一张表，用来缓存用以沟通的数据表，其中相关的rpc连接要用的时候再去连
}

// ToDo:直接返回一个查询体，主要需要解析出来 ①查询的table名称 ②执行的操作类型
func parser(input string) (regionRPC.SQLRst, bool) {
	rst := regionRPC.SQLRst{}
	word := strings.Split(input, " ")
	// [检查]是否以分号结尾
	if input[len(input)-1] != ';' {
		fmt.Println("SQL语句没有以分号结尾!")
		return rst, false
	}
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
func runSQL(input string) {
	// 1. 得到解析后的SQL内容
	sqlRst, ok := parser(input)
	if ok {
		// 2. 尝试在本地查找相关的表
		for _, table := range global.TableCache {
			if table.Name == sqlRst.Table {
				fmt.Println("table '" + table.Name + "' 在本地有缓存")
			}
		}
	} else {
		fmt.Println("错误的SQL语句")
	}
}
