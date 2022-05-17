package etcd

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

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

func GetHostAddress() string {
	return "127.0.0.1"
}
