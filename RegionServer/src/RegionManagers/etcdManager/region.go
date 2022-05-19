package etcd

import (
	"context"
	"fmt"
	"log"
	config "region/utils/ConfigSystem"
	mylog "region/utils/LogSystem"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 获取一个etcd的连接
func Init() *clientv3.Client {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{config.Configs.Etcd_ip + ":" + config.Configs.Etcd_port},
		DialTimeout: 5 * time.Second,
	})
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

	_, err = client.Put(ctx, config.Configs.Etcd_region_register_catalog+"/"+GetHostAddress(), "", clientv3.WithLease(leaseGrant.ID))
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
		log := mylog.NewNormalLog("服务器" + GetHostAddress() + "尝试续约")
		log.LogType = "INFO"
		log.LogGen(mylog.LogInputChan)
		<-keepaliveResponseChan
		// fmt.Println("ttl:", ka.TTL)
	}

}

// 获取master的相关信息(返回 ip+port)
func GetMasterAddress(client *clientv3.Client) string {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	getResponse, err := client.Get(ctx, config.Configs.Etcd_master_address)
	if err != nil {
		log.Printf("etcd GET error,%v\n", err)
		return ""
	}

	for _, kv := range getResponse.Kvs {
		fmt.Printf("%s=%s\n", kv.Key, kv.Value)
	}
	return string(getResponse.Kvs[0].Value)
}

// 返回自己的IP地址
func GetHostAddress() string {
	return "127.0.0.1"
}

//=============主从复制=============
// 方法：分区服务器提交自己的版本号
func AddVersion(client *clientv3.Client, tableName string, version string) error {
	return nil
}

// 方法：分区服务器获得表的主副本的IP地址+端口号
func GetMaster(client *clientv3.Client, tableName string) string {
	return ""
}

// 方法：分区服务器获得表的syncCopys的IP地址+端口号
func GetSyncCopys(client *clientv3.Client, tableName string) string {
	return ""
}

// 方法：监听本地的主副本的目录，一旦有别的分区服务器加入，则进行一些操作
