package rpc

import "net/rpc"

// 最终暴露给客户端的服务主体
type ReportServiceClient struct {
	*rpc.Client
}

// 接口定义
type ReportTableServiceInterface = interface {
	ReportTable(request []LocalTable, reply *string) error
}

// 接口断言
var _ ReportTableServiceInterface = (*ReportServiceClient)(nil)

// 客户端服务生成函数
func DialReportService(network, address string) (*ReportServiceClient, error) {
	c, err := rpc.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &ReportServiceClient{Client: c}, nil
}

// 最终暴露给用户的服务主体实现的方法
func (p *ReportServiceClient) ReportTable(request []LocalTable, reply *string) error {
	return p.Client.Call(ReportServiceName+".ReportTable", request, reply)
}

//以下是方便服务端开发的代码

const ReportServiceName = "ReportService" // 服务的名字，使用包前缀

// 方便服务端绑定服务的函数，同时检测服务端是否真正实现了这个接口
func RegisterReportService(svc ReportTableServiceInterface) error {
	return rpc.RegisterName(ReportServiceName, svc)
}
