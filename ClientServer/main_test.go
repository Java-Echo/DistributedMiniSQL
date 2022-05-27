package main

import (
	masterRPC "client/rpcManager/master"
	regionRPC "client/rpcManager/region"
	config "client/utils/ConfigSystem"
	mylog "client/utils/LogSystem"
	"client/utils/global"
	"fmt"
	"log"
	"strconv"
	"testing"
)

func TestMain(m *testing.M) {
	mylog.LogInputChan = mylog.LogStart()
	config.BuildConfig()
	var err error
	masterRPC.RpcM2R, err = masterRPC.DialService("tcp", config.Configs.Master_ip+":"+config.Configs.Master_port)
	if err != nil {
		log.Fatal(err)
	}
	global.TableCache = make(map[string]global.TableMeta)
	// make一下哦
	// global.TableCache
	// client, _ = masterRpc.DialService("tcp", "localhost:"+config.Configs.Master_port)
	m.Run()

}

func Test_parser(t *testing.T) {
	testSQL := []string{
		"select * from student2 where (((213!=id or 70<score) and 'name123'<name) or 99.0!=score) and name>'name9998';",
		"insert into student2 values(10866666666,'王振阳',100);",
		"update student2 set name='陈旭征' where name='王振阳';",
		"delete from student2 where name='name97996';",
		"drop table student2;",
		"create table ttt;",
	}
	for i, sql := range testSQL {
		s, _ := parser(sql)
		fmt.Println("-------测试第" + strconv.Itoa(i) + "条语句-------")
		fmt.Println("SQL语句为:" + s.SQL)
		fmt.Println("SQL的类型为:" + s.SQLtype)
		fmt.Println("SQL执行的表为:" + s.Table)
	}
	// s := parser("select * from aaa;")
	// fmt.Println(s.SQL)
	// fmt.Println(s.SQLtype)
	// fmt.Println(s.Table)
	t.Error("终止")
}

func Test_runOnRegion(t *testing.T) {
	sql := regionRPC.SQLRst{
		SQLtype: "select",
		SQL:     "select * from ttt;",
		Table:   "ttt",
	}
	res, err := runOnRegion(sql, "192.168.31.68")
	if err != nil {
		log.Fatal("RunOnRegion error:", err)
	}
	fmt.Println("sql返回的结果为:" + res)
	t.Error("终止")
}

// 尚未完成测试
func Test_chooseRegionAndRun(t *testing.T) {
	// SQL语句
	sql := regionRPC.SQLRst{
		SQLtype: "select",
		SQL:     "select * from ttt;",
		Table:   "ttt",
	}

	// tableMeta
	meta := global.TableMeta{
		Name: "ttt",
		Master: global.RegionInfo{
			IP: "192.168.31.68",
		},
		Sync_slave: global.RegionInfo{
			// IP: "192.168.31.68",
		},
		Slaves: []global.RegionInfo{
			// global.RegionInfo{IP: "192.168.31.68"},
			// global.RegionInfo{IP: "192.168.31.68"},
		},
	}

	// 测试运行
	res, err := chooseRegionAndRun(sql, meta)
	if err != nil {
		log.Fatal("runSQL error:", err)
	}
	fmt.Println("sql语句查询的结果为:" + res)

	t.Error("终止")
}

// 尚未完成测试
func Test_runSQL(t *testing.T) {
	// sql := "select * from cyy;"
	sql := "insert into hyh values(6);"
	res, err := runSQL(sql)
	if err != nil {
		log.Fatal("runSQL error:", err)
	}
	fmt.Println("sql语句查询的结果为:" + res)
	t.Error("终止")
}
