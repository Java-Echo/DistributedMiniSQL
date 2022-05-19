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
