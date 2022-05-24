package master

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
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

// 工具函数：得到路径的最后一个字段
func util_getLastKey(path string) string {
	keys := strings.Split(path, "/")
	return keys[len(keys)-1]
}

// ToDo:合理安排这张全局的表的位置
var RegionMap = make(map[string]string)

// 进行相关的配置
func Init() *clientv3.Client {
	global.HostIP = GetHostAddress()
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{config.Configs.Etcd_ip + ":" + config.Configs.Etcd_port},
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
