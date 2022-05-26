package etcd

import (
	"context"
	"fmt"
	"log"
	"net"
	config "region/utils/ConfigSystem"
	mylog "region/utils/LogSystem"
	"region/utils/global"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 获取一个etcd的连接
func Init() *clientv3.Client {
	global.HostIP = GetHostAddress()
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{config.Configs.Etcd_ip + ":" + config.Configs.Etcd_port},
		DialTimeout: 5 * time.Second,
	})
	log_ := mylog.NewNormalLog("成功连入etcd")
	log_.LogType = "INFO"
	log_.LogGen(mylog.LogInputChan)
	if err != nil {
		log.Fatalln(err)
	}

	return client
}

//=============服务发现=============

// 向etcd注册自己
func ServiceRegister(client *clientv3.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 尝试注册一个新的租约
	// 获取一个租约 有效期为5秒
	leaseGrant, err := client.Grant(ctx, 5)
	if err != nil {
		log.Printf("put error %v", err)
		return
	}

	_, err = client.Put(ctx, config.Configs.Etcd_region_register_catalog+"/"+global.HostIP, "", clientv3.WithLease(leaseGrant.ID))
	if err != nil {
		log.Printf("put error %v", err)
		return
	}

	// 发送心跳包，不断续约
	keepaliveResponseChan, err := client.KeepAlive(ctx, leaseGrant.ID)
	if err != nil {
		log.Printf("KeepAlive error %v", err)
		return
	}

	if err != nil {
		log.Printf("KeepAlive error %v", err)
		return
	}

	for {
		// log := mylog.NewNormalLog("服务器" + global.HostIP + "尝试续约")
		// log.LogType = "INFO"
		// log.LogGen(mylog.LogInputChan)
		time.Sleep(1 * time.Second)
		<-keepaliveResponseChan
		// fmt.Println("ttl:", ka.TTL)
	}
}

// 获取master的相关信息(返回 ip+port)
func GetMasterIP(client *clientv3.Client) string {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	getResponse, err := client.Get(ctx, config.Configs.Etcd_master_address)
	if err != nil {
		log.Printf("etcd GET error,%v\n", err)
		return ""
	}

	// for _, kv := range getResponse.Kvs {
	// 	fmt.Printf("%s=%s\n", kv.Key, kv.Value)
	// }
	return string(getResponse.Kvs[0].Value)
}

// 返回自己的IP地址
func GetHostAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	fmt.Println("怎么出来了？")
	return "127.0.0.1"
}

func GetWatcher(client *clientv3.Client, catalog string) *clientv3.WatchChan {
	watchChan := client.Watch(context.Background(), catalog, clientv3.WithPrefix())
	log := mylog.NewNormalLog("开启对目录" + catalog + "的监听")
	log.LogGen(mylog.LogInputChan)
	return &watchChan
}

// 工具函数：得到路径的最后一个字段
func util_getLastKey(path string) string {
	keys := strings.Split(path, "/")
	return keys[len(keys)-1]
}

//=============主从复制=============
// 方法：分区服务器提交自己的版本号
func AddVersion(tableName string, version string) error {
	return nil
}

// 方法：分区服务器获得表的主副本的IP地址+端口号
func GetMaster(tableName string) string {
	return ""
}

// 方法：分区服务器获得表的syncCopy的IP地址+端口号
func GetSyncSlave(tableName string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	getResponse, err := global.Region.Get(ctx, config.Configs.Etcd_table_catalog+"/"+tableName+"/sync_slave/", clientv3.WithPrefix())
	if err != nil {
		log.Printf("etcd GET error,%v\n", err)
		return ""
	}
	for _, kv := range getResponse.Kvs {
		return util_getLastKey(string(kv.Key)) + ":" + string(kv.Value)
	}
	return ""
}

// 方法：分区服务器获得表的异步copy的IP地址+端口号
func GetSlaves(tableName string) []string {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	getResponse, err := global.Region.Get(ctx, config.Configs.Etcd_table_catalog+"/"+tableName+"/slave/", clientv3.WithPrefix())
	if err != nil {
		log.Printf("etcd GET error,%v\n", err)
		return nil
	}
	var res []string
	for _, kv := range getResponse.Kvs {
		res = append(res, util_getLastKey(string(kv.Key))+":"+string(kv.Value))
	}
	return res
}

// 方法：监听本地的主副本的目录，一旦有别的分区服务器加入，则进行一些操作
