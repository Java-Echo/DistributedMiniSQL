package regionRpc

/*-------------Region--------------------*/
type SQLRst struct {
	SQLtype string // SQL语句的类型
	Table   string // SQL语句具体查询的表
	SQL     string // 具体的SQL语句
}

type SQLRes struct {
	State  string // 查询结果的状态(成功、失败等)
	Result string // 最终SQL的返回结果
}
