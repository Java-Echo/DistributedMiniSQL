package etcd

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 获取一个etcd的连接
func Init() *clientv3.Client {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return client
}

//=============服务发现=============

// 向etcd注册自己
func ServiceRegister(client *clientv3.Client, catalog string) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 得到当前序号
	// ToDo:这里的catalog理应是在配置文件中的一个内容
	getResponse, err := client.Get(ctx, catalog, clientv3.WithPrefix())
	if err != nil {
		log.Printf("etcd put error,%v\n", err)
		return
	}

	newIndex := 1
	if len(getResponse.Kvs) != 0 {
		lastName := getResponse.Kvs[len(getResponse.Kvs)-1].Key
		lastIndex, _ := strconv.Atoi(string(lastName[len(catalog):]))
		newIndex = lastIndex + 1
	}

	// 尝试注册一个新的租约
	// 获取一个租约 有效期为5秒
	leaseGrant, err := client.Grant(ctx, 5)
	if err != nil {
		log.Printf("put error %v", err)
		return
	}

	fmt.Println("当前服务器的地址为：" + GetHostAddress())
	fmt.Println("获得的租约服务器名为:" + catalog + "/region_server" + strconv.Itoa(newIndex))
	_, err = client.Put(ctx, catalog+"/region_server"+strconv.Itoa(newIndex), GetHostAddress(), clientv3.WithLease(leaseGrant.ID))
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
		<-keepaliveResponseChan
		// fmt.Println("ttl:", ka.TTL)
	}

}

// 获取master的相关信息(返回 ip+port)
func GetMasterConfig(client *clientv3.Client, key string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	getResponse, err := client.Get(ctx, key)
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
