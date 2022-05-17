package rpc

/*--------------ReportService--------------------*/

type LocalTable struct {
	Name string
	IP   string
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
