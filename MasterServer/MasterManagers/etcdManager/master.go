package master

import (
	"context"
	"fmt"
	"log"
	"time"

	mylog "master/utils/LogSystem"
	gloabl "master/utils/global"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// ToDo:合理安排这张全局的表的位置
var TableMap = make(map[string]string)

// 进行相关的配置
func Init() *clientv3.Client {
	fmt.Println("尝试初始化etcd连接")
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("etcd成功连接")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 上传配置信息
	catalog := "/config"
	_, err = client.Put(ctx, catalog+"/masterAddress", "测试") // ToDo：得到master需要配置的地址
	// 其他配置信息
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("成功启动master服务器")
	return client
}

//=============服务发现=============

// ToDo:根据监控到的改变数据进行本地Region服务器的调整
func RegisterWatcher(client *clientv3.Client, catalog string) {
	watchChan := client.Watch(context.Background(), catalog, clientv3.WithPrefix())
	log := mylog.NewNormalLog("master节点开启节点目录监听")
	log.LogGen(mylog.LogInputChan)

	for watchResponse := range watchChan {
		for _, event := range watchResponse.Events {
			fmt.Printf("Type:%s,Key:%s,Value:%s\n", event.Type, event.Kv.Key, event.Kv.Value)
			if event.Type == 0 {
				// 为新加入的节点添加元信息
				newMeta := gloabl.RegionMeta{}
				// ToDo:进一步完善相关的信息
				gloabl.TableMap[string(event.Kv.Key)] = newMeta
				// 记录日志
				log := mylog.NewNormalLog("服务器 " + string(event.Kv.Key) + " 尝试建立连接")
				log.LogGen(mylog.LogInputChan)
			} else if event.Type == 1 {
				// 删除新加入节点的元信息
				delete(gloabl.TableMap, string(event.Kv.Key))
				// 记录日志
				log := mylog.NewNormalLog("服务器 " + string(event.Kv.Key) + " 失去连接")
				log.LogGen(mylog.LogInputChan)
			}
		}
	}
}

//=============主从复制=============

// 方法：主服务器为一个从副本建立/删除数据表下的注册
func CreateSlave(client *clientv3.Client, tableName string, ip string) error {
	return nil
}

func DeleteSlave(client *clientv3.Client, tableName string, ip string) error {
	return nil
}

// 方法：主服务器为master建立/删除注册
func CreateMaster(client *clientv3.Client, tableName string, ip string) error {
	return nil
}

func DeleteMaster(client *clientv3.Client, tableName string, ip string) error {
	return nil
}

// 方法：主服务器为syncCopys建立/删除注册
func CreateSyncCopys(client *clientv3.Client, tableName string, ip string) error {
	return nil
}

func DeleteSyncCopys(client *clientv3.Client, tableName string, ip string) error {
	return nil
}
