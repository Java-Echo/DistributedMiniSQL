package rpc

// 用于记录rpc中特定的消息传输格式

/*--------------ReportService--------------------*/
type LocalTable struct {
	Name string
	IP   string
	Port string
}

type ValidTable struct {
	Name  string
	Level string //指明是主副本还是(异步)从副本
}

type ReportTableRes struct {
	Tables []ValidTable
}

/*--------------CliService--------------------*/
type Region struct {
	IP string
}

type TableInfo struct {
	Name       string
	Master     Region   // 主副本所在的节点
	Sync_slave Region   // 同步从副本所在的节点
	Slaves     []Region // 异步从副本所在的节点
}
