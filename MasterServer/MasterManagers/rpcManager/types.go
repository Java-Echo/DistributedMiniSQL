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
