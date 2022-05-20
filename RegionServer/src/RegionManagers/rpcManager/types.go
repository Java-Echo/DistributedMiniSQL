package rpc

/*--------------ReportService--------------------*/

type LocalTable struct {
	Name  string
	IP    string
	Port  string
	Level string
}

type ValidTable struct {
	Name  string
	Level string //指明是主副本还是(异步)从副本
}

type ReportTableRes struct {
	Tables []ValidTable
}

/*-------------GossipService--------------------*/
type FetchLogRst struct {
	TableName string
	Version   string // 本地的版本号
}

type FetchLogRes struct {
	TableName string
	Version   string
	Log       []string
}

type SyncProbeRst struct {
	TableName string
	Version   string
}

type SyncProbeRes struct {
}

/*-------------CliService--------------------*/
type SQLRst struct {
	SQLtype string // SQL语句的类型
	Table   string // SQL语句具体查询的表
	SQL     string // 具体的SQL语句
}

type SQLRes struct {
	state  string // 查询结果的状态(成功、失败等)
	result string // 最终SQL的返回结果
}
