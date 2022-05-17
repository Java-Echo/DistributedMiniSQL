package master

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// ToDo:合理安排这张全局的表的位置
var TableMap = make(map[string]string)

// 进行相关的配置
func Init() *clientv3.Client {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalln(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 上传配置信息
	catalog := "config"
	_, err = client.Put(ctx, catalog+"/masterAddress", "这里应该填写地址") // ToDo：得到master需要配置的地址
	// 其他配置信息
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("成功启动master服务器")
	return client
}

// ToDo:根据监控到的改变数据进行本地Region服务器的调整
func RegisterWatcher(catalog string, client *clientv3.Client) {
	watchChan := client.Watch(context.Background(), catalog, clientv3.WithPrefix())
	fmt.Println("正在监听" + catalog)

	for watchResponse := range watchChan {
		for _, event := range watchResponse.Events {
			fmt.Printf("Type:%s,Key:%s,Value:%s\n", event.Type, event.Kv.Key, event.Kv.Value)
		}
	}

}
