package rpc

import "net/rpc"

// 最终暴露给客户端的服务主体
type GossipServiceClient struct {
	*rpc.Client
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

//以下是方便服务端开发的代码

const GossipServiceName = "GossipService" // 服务的名字，使用包前缀

// 方便服务端绑定服务的函数，同时检测服务端是否真正实现了这个接口
func RegisterGossipService(svc GossipServiceInterface) error {
	return rpc.RegisterName(GossipServiceName, svc)
}

// ========以下是真正暴露出来的接口========
// 接口定义
type GossipServiceInterface = interface {
	PassTable(request PassTableRst, reply *PassTableRes) error
	PassLog(request PassLogRst, reply *PassLogRes) error
	// ToDo:可以定义更多的函数
}

// 传递一条日志
func (p *GossipServiceClient) PassLog(request PassLogRst, reply *PassLogRes) error {
	return p.Client.Call(GossipServiceName+".PassLog", request, reply)
}

// 传递一个表的文件
func (p *GossipServiceClient) PassTable(request PassTableRst, reply *PassTableRes) error {
	return p.Client.Call(GossipServiceName+".PassTable", request, reply)
}
