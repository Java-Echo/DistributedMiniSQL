package rpc

import (
	"fmt"
	"net/rpc"
	mylog "region/utils/LogSystem"
)

// 最终暴露给客户端的服务主体
type CliServiceClient struct {
	*rpc.Client
}

// 接口断言
var _ CliServiceInterface = (*CliServiceClient)(nil)

// 客户端服务生成函数
func DialService(network, address string) (*CliServiceClient, error) {
	c, err := rpc.Dial(network, address)
	if err != nil {
		fmt.Println("DialService没找到服务")
		return nil, err
	}
	log_ := mylog.NewNormalLog("成功注册DialService")
	log_.LogType = "INFO"
	log_.LogGen(mylog.LogInputChan)
	return &CliServiceClient{Client: c}, nil
}

//以下是方便服务端开发的代码

const ServiceName = "Service" // 服务的名字，使用包前缀

// 方便服务端绑定服务的函数，同时检测服务端是否真正实现了这个接口
func RegisterCliService(svc CliServiceInterface) error {
	return rpc.RegisterName(ServiceName, svc)
}

// ========以下是真正暴露出来的接口========
// 接口定义
type CliServiceInterface = interface {
	Hello(request string, reply *string) error
	SQL(request SQLRst, reply *SQLRes) error
}

// 实现方法1
func (p *CliServiceClient) Hello(request string, reply *string) error {
	return p.Client.Call(ServiceName+".Hello", request, reply)
}

// 方法：SQL的执行
func (p *CliServiceClient) SQL(request SQLRst, reply *SQLRes) error {
	return p.Client.Call(ServiceName+".SQL", request, reply)
}
