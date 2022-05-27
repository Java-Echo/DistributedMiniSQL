package master

import (
	"context"
	"fmt"
	"log"
	config "master/utils/ConfigSystem"
	mylog "master/utils/LogSystem"
	"master/utils/global"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// ToDo:根据监控到的改变数据进行本地Region服务器的调整
func RegisterWatcherWithWorker(client *clientv3.Client, catalog string, worker WatchWorkerInterface) {
	watchChan := client.Watch(context.Background(), catalog, clientv3.WithPrefix())
	log := mylog.NewNormalLog("master节点开启对集群节点目录'" + catalog + "'的监听")
	log.LogGen(mylog.LogInputChan)

	for watchResponse := range watchChan {
		for _, event := range watchResponse.Events {
			// fmt.Printf("Type:%s,Key:%s,Value:%s\n", event.Type, event.Kv.Key, event.Kv.Value)
			if event.Type == clientv3.EventTypePut {
				worker.OnPut(event)
			} else if event.Type == clientv3.EventTypeDelete {
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
	ip := util_getLastKey(string(event.Kv.Key))
	table, ok := global.RegionMap[ip]

	if ok {
		// 服务器是在一定时间里面宕机重启
		if table.State == global.Stop {
			log := mylog.NewNormalLog("服务器 " + ip + " 正在尝试宕机重启")
			log.LogGen(mylog.LogInputChan)
			table.State = global.Working

			global.PrintRegionMap(1)
			global.PrintTableMap(1)
		}
		// 如果不是的话，那就是服务器在尝试续约
	} else {
		// 服务器是新加入的，为新加入的节点添加元信息
		log := mylog.NewNormalLog("服务器 " + ip + " 尝试建立连接")
		log.LogGen(mylog.LogInputChan)
		newMeta := &global.RegionMeta{}
		newMeta.IP = ip
		newMeta.State = global.Working
		// 将该节点加入到全局的表中
		global.RegionMap[ip] = newMeta

		global.PrintRegionMap(1)
		global.PrintTableMap(1)
	}
}

func (p *RegionRegisterWorker) OnDelete(event *clientv3.Event) {
	// fmt.Printf("Type:%s,Key:%s,Value:%s\n", event.Type, event.Kv.Key, event.Kv.Value)
	ip := util_getLastKey(string(event.Kv.Key))

	log_ := mylog.NewNormalLog("服务器 " + ip + " 失去连接")
	log_.LogGen(mylog.LogInputChan)

	// ToDo:修改全局表中的相关信息，这里的逻辑需要完善
	preMeta := global.RegionMap[ip]
	preMeta.State = global.Stop
	// 从本地的table信息映射表+etcd目录中删除相关的节点
	for _, table := range global.TableMap {
		if table.MasterRegion == ip {
			DeleteMaster(table.Name, ip)
			// ToDo:这里需要启用从副本来进行容错容灾
			table.MasterRegion = ""
			if table.SyncRegion != "" {
				// 将同步从副本提升上来(尚未测试)
				table.MasterRegion = table.SyncRegion
				table.SyncRegion = ""
				CreateMaster(table.Name, table.MasterRegion)
				DeleteSyncSlave(table.Name, table.SyncRegion)
			} else {
				fmt.Println("主副本没了,也没有从副本,寄了")
			}
		} else if table.SyncRegion == ip {
			DeleteSyncSlave(table.Name, ip)
			table.SyncRegion = ""
		}
		for i, slave := range table.CopyRegions {
			if slave == ip {
				DeleteSlave(table.Name, ip)
				table.CopyRegions = append(table.CopyRegions[:i], table.CopyRegions[i+1:]...)
			}
		}
	}

	// 将该region服务器加入到暂存区
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Configs.Etcd_region_stepout_time)*time.Second)
	defer cancel()
	// 获取一个租约
	leaseGrant, err := global.Master.Grant(ctx, config.Configs.Etcd_region_stepout_time)
	if err != nil {
		log.Printf("put error %v", err)
		return
	}

	_, err = global.Master.Put(ctx, config.Configs.Etcd_region_stepout_catalog+"/"+ip, "", clientv3.WithLease(leaseGrant.ID))
	if err != nil {
		log.Printf("put error %v", err)
		return
	}

	log_ = mylog.NewNormalLog("服务器 " + ip + " 加入到了暂存区")
	log_.LogGen(mylog.LogInputChan)

	global.PrintRegionMap(1)
	global.PrintTableMap(1)
}

var _ WatchWorkerInterface = (*RegionRegisterWorker)(nil)

/*---------------对于region暂时失去连接的目录的监听---------------*/
type RegionStepOutWorker struct {
}

// ToDo:这里的逻辑是不完善的
func (p *RegionStepOutWorker) OnPut(event *clientv3.Event) {
	// fmt.Printf("Type:%s,Key:%s,Value:%s\n", event.Type, event.Kv.Key, event.Kv.Value)
	ip := util_getLastKey(string(event.Kv.Key))

	log := mylog.NewNormalLog("服务器 " + ip + " 进入“暂时失联”状态")
	log.LogGen(mylog.LogInputChan)
}

func (p *RegionStepOutWorker) OnDelete(event *clientv3.Event) {
	// fmt.Printf("Type:%s,Key:%s,Value:%s\n", event.Type, event.Kv.Key, event.Kv.Value)
	ip := util_getLastKey(string(event.Kv.Key))
	log := mylog.NewNormalLog("服务器 " + ip + " 离开“暂时失联”状态")
	log.LogGen(mylog.LogInputChan)

	regionMeta := global.RegionMap[ip]
	if regionMeta.State == global.Stop {
		log := mylog.NewNormalLog("服务器 " + ip + " 完全失去联系, 正在尝试清除它的一切")
		log.LogGen(mylog.LogInputChan)
		// 从全局的region服务器表中删除(在这之前我们已经将其从全部的表中删除了)
		delete(global.RegionMap, regionMeta.IP)
	} else if regionMeta.State == global.Working {
		log := mylog.NewNormalLog("服务器 " + ip + " 宕机重启成功")
		log.LogGen(mylog.LogInputChan)
	}

	global.PrintRegionMap(1)
	global.PrintTableMap(1)
}

var _ WatchWorkerInterface = (*RegionStepOutWorker)(nil)
