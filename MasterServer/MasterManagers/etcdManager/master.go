package master

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	config "master/utils/ConfigSystem"
	mylog "master/utils/LogSystem"
	"master/utils/global"

	clientv3 "go.etcd.io/etcd/client/v3"
)

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

// ToDo:合理安排这张全局的表的位置
var RegionMap = make(map[string]string)

// 进行相关的配置
func Init() *clientv3.Client {
	global.HostIP = GetHostAddress()
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log_ := mylog.NewNormalLog("成功连入etcd")
	log_.LogType = "INFO"
	log_.LogGen(mylog.LogInputChan)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 上传配置信息
	catalog := "/config"
	_, err = client.Put(ctx, catalog+"/masterAddress", global.HostIP) // ToDo：得到master需要配置的地址
	// 其他配置信息
	if err != nil {
		log.Fatalln(err)
	}
	log_ = mylog.NewNormalLog("成功上传配置信息")
	log_.LogType = "INFO"
	log_.LogGen(mylog.LogInputChan)
	return client
}

//=============服务发现=============

// ToDo:根据监控到的改变数据进行本地Region服务器的调整
func RegisterWatcher(client *clientv3.Client, catalog string) {
	watchChan := client.Watch(context.Background(), catalog, clientv3.WithPrefix())
	log := mylog.NewNormalLog("master节点开启对集群节点目录监听")
	log.LogGen(mylog.LogInputChan)

	for watchResponse := range watchChan {
		for _, event := range watchResponse.Events {
			fmt.Printf("Type:%s,Key:%s,Value:%s\n", event.Type, event.Kv.Key, event.Kv.Value)
			if event.Type == 0 {
				// 为新加入的节点添加元信息
				newMeta := global.RegionMeta{}
				// ToDo:进一步完善相关的信息
				global.RegionMap[string(event.Kv.Key)] = newMeta
				// 记录日志
				log := mylog.NewNormalLog("服务器 " + string(event.Kv.Key) + " 尝试建立连接")
				log.LogGen(mylog.LogInputChan)
			} else if event.Type == 1 {
				// 删除新加入节点的元信息
				delete(global.RegionMap, string(event.Kv.Key))
				// 记录日志
				log := mylog.NewNormalLog("服务器 " + string(event.Kv.Key) + " 失去连接")
				log.LogGen(mylog.LogInputChan)
			}
		}
	}
}

//=============主从复制=============

// 方法：主服务器为一个从副本建立/删除数据表下的注册
func CreateSlave(tableName string, ip string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := global.Master.Put(ctx, config.Configs.Etcd_table_catalog+"/"+tableName+"/"+"slave/"+ip, "")
	if err != nil {
		log.Fatalln(err)
		return err
	}
	log := mylog.NewNormalLog("将节点" + ip + "加入到了表" + tableName + "的slave副本下")
	log.LogGen(mylog.LogInputChan)

	return nil
}

func DeleteSlave(table global.TableMeta, ip string) error {
	return nil
}

// 方法：主服务器为master建立/删除注册
func CreateMaster(tableName string, ip string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := global.Master.Put(ctx, config.Configs.Etcd_table_catalog+"/"+tableName+"/"+"master/"+ip, "")
	if err != nil {
		log.Fatalln(err)
		return err
	}
	log := mylog.NewNormalLog("将节点" + ip + "加入到了表" + tableName + "的master副本下")
	log.LogGen(mylog.LogInputChan)

	return nil
}

func DeleteMaster(table global.TableMeta, ip string) error {

	return nil
}

// 方法：主服务器为syncCopys建立/删除注册
func CreateSyncSlave(tableName string, ip string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := global.Master.Put(ctx, config.Configs.Etcd_table_catalog+"/"+tableName+"/"+"sync_slave/"+ip, "")
	if err != nil {
		log.Fatalln(err)
		return err
	}
	log := mylog.NewNormalLog("将节点" + ip + "加入到了表" + tableName + "的sync_slave副本下")
	log.LogGen(mylog.LogInputChan)

	return nil
}

func DeleteSyncSlave(table global.TableMeta, ip string, port string) error {
	return nil
}
