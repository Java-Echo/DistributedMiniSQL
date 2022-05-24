package master

import (
	"context"
	"fmt"
	mylog "master/utils/LogSystem"
	"master/utils/global"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func RegisterWatcher(client *clientv3.Client, catalog string) {
	watchChan := client.Watch(context.Background(), catalog, clientv3.WithPrefix())
	log := mylog.NewNormalLog("master节点开启对集群节点目录'" + catalog + "'的监听")
	log.LogGen(mylog.LogInputChan)

	for watchResponse := range watchChan {
		for _, event := range watchResponse.Events {
			fmt.Printf("Type:%s,Key:%s,Value:%s\n", event.Type, event.Kv.Key, event.Kv.Value)
			if event.Type == clientv3.EventTypePut {
				// 为新加入的节点添加元信息
				newMeta := global.RegionMeta{}
				// ToDo:进一步完善相关的信息
				global.RegionMap[string(event.Kv.Key)] = newMeta
				// 记录日志
				log := mylog.NewNormalLog("服务器 " + string(event.Kv.Key) + " 尝试建立连接")
				log.LogGen(mylog.LogInputChan)
			} else if event.Type == clientv3.EventTypeDelete {
				// 记录日志
				log := mylog.NewNormalLog("服务器 " + string(event.Kv.Key) + " 失去连接")
				log.LogGen(mylog.LogInputChan)
			}
		}
	}
}

// ToDo:根据监控到的改变数据进行本地Region服务器的调整
func RegisterWatcherWithWorker(client *clientv3.Client, catalog string, worker WatchWorkerInterface) {
	watchChan := client.Watch(context.Background(), catalog, clientv3.WithPrefix())
	log := mylog.NewNormalLog("master节点开启对集群节点目录'" + catalog + "'的监听")
	log.LogGen(mylog.LogInputChan)

	for watchResponse := range watchChan {
		for _, event := range watchResponse.Events {
			fmt.Printf("Type:%s,Key:%s,Value:%s\n", event.Type, event.Kv.Key, event.Kv.Value)
			if event.Type == 0 {
				worker.OnPut(event)
			} else if event.Type == 1 {
				worker.OnDelete(event)
			}
		}
	}
}

type WatchWorkerInterface = interface {
	OnPut(event *clientv3.Event)
	OnDelete(event *clientv3.Event)
}

/*---------------对于region注册目录(etcd_region_register_catalog)的监听---------------*/
type RegionRegisterWorker struct {
}

func (p *RegionRegisterWorker) OnPut(event *clientv3.Event) {
	log := mylog.NewNormalLog("服务器 " + string(event.Kv.Key) + " 尝试建立连接")
	log.LogGen(mylog.LogInputChan)
}
func (p *RegionRegisterWorker) OnDelete(event *clientv3.Event) {
	log := mylog.NewNormalLog("服务器 " + string(event.Kv.Key) + " 失去连接")
	log.LogGen(mylog.LogInputChan)
}

var _ WatchWorkerInterface = (*RegionRegisterWorker)(nil)

/*---------------对于region暂时失去连接的目录的监听---------------*/
type RegionStepOutWorker struct {
}

func (p *RegionStepOutWorker) OnPut(event *clientv3.Event) {
	log := mylog.NewNormalLog("服务器 " + string(event.Kv.Key) + " 进入“暂时失联”状态")
	log.LogGen(mylog.LogInputChan)
}
func (p *RegionStepOutWorker) OnDelete(event *clientv3.Event) {
	log := mylog.NewNormalLog("服务器 " + string(event.Kv.Key) + " 离开“暂时失联”状态")
	log.LogGen(mylog.LogInputChan)
}

var _ WatchWorkerInterface = (*RegionStepOutWorker)(nil)
