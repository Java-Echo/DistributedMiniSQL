package rpc

import "net/rpc"

// 最终暴露给客户端的服务主体
type GossipServiceClient struct {
	*rpc.Client
}

// 接口定义
type GossipServiceInterface = interface {
	FetchLog(request FetchLogRst, reply *FetchLogRes) error    // 其他region服务器来询问特定表最新的log
	SyncProbe(request SyncProbeRst, reply *SyncProbeRes) error // 其他region服务器来询问是否完成了特定表特定版本的同步
	// ToDo:可以定义更多的函数
}

// 接口断言
var _ GossipServiceInterface = (*GossipServiceClient)(nil)

// 客户端服务生成函数
func DialGossipService(network, address string) (*GossipServiceClient, error) {
	c, err := rpc.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &GossipServiceClient{Client: c}, nil
}

// 实现方法1
func (p *GossipServiceClient) FetchLog(request FetchLogRst, reply *FetchLogRes) error {
	return p.Client.Call(GossipServiceName+".FetchLog", request, reply)
}

// 实现方法2
func (p *GossipServiceClient) SyncProbe(request SyncProbeRst, reply *SyncProbeRes) error {
	return p.Client.Call(GossipServiceName+".SyncProbe", request, reply)
}

//以下是方便服务端开发的代码

const GossipServiceName = "GossipService" // 服务的名字，使用包前缀

// 方便服务端绑定服务的函数，同时检测服务端是否真正实现了这个接口
func RegisterGossipService(svc GossipServiceInterface) error {
	return rpc.RegisterName(GossipServiceName, svc)
}
